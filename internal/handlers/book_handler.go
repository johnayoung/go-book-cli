package handlers

import (
	"fmt"
	"go-book-ai/internal/agents"
	"go-book-ai/internal/errors"
	"go-book-ai/internal/file"
	"go-book-ai/internal/logger"
	"go-book-ai/internal/models"
	"go-book-ai/internal/outline"
	"go-book-ai/internal/state"
	"go-book-ai/internal/utils"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// BookCommandHandler handles book-related commands.
type BookCommandHandler struct {
	WritingAgent   agents.WritingAgent
	ReviewingAgent agents.ReviewingAgent
	FileManager    *file.FileManager
	ErrorHandler   *errors.ErrorHandler
	Logger         logger.Logger
}

// NewBookCommandHandler returns a new BookCommandHandler.
func NewBookCommandHandler(writingAgent agents.WritingAgent, reviewingAgent agents.ReviewingAgent, fm *file.FileManager, eh *errors.ErrorHandler, lg logger.Logger) *BookCommandHandler {
	return &BookCommandHandler{WritingAgent: writingAgent, ReviewingAgent: reviewingAgent, FileManager: fm, ErrorHandler: eh, Logger: lg}
}

func (h *BookCommandHandler) handleError(msg string, err error) error {
	h.ErrorHandler.LogError(fmt.Errorf("%s: %v", msg, err))
	return fmt.Errorf("%s: %v", msg, err)
}

func (h *BookCommandHandler) ProcessBook(topic string) error {
	folderName := utils.CleanName(topic)
	bookPath := filepath.Join("books", folderName)
	h.Logger.Info(fmt.Sprintf("Processing book folder: %s", bookPath))

	var history models.ConversationHistory
	historyPath := filepath.Join(bookPath, "history.json")

	// Create the book directory if it does not exist
	err := os.MkdirAll(bookPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create book directory: %w", err)
	}

	// Attempt to load conversation history
	err = h.FileManager.LoadHistoryFromFile(historyPath, &history)
	if err != nil {
		if os.IsNotExist(err) {
			history = models.StartNewConversation(topic)
		} else {
			return fmt.Errorf("failed to load conversation history: %w", err)
		}
	}

	bookState, err := state.LoadState(bookPath)
	if err != nil {
		if os.IsNotExist(err) {
			bookState = state.NewState()
		} else {
			return fmt.Errorf("failed to load state: %v", err)
		}
	}

	// Check if the outline has already been generated
	if !bookState.OutlineGenerated {
		err := h.generateBookOutline(topic, bookPath, bookState, &history, historyPath)
		if err != nil {
			return err
		}
	} else {
		h.Logger.Info("Outline already generated, skipping outline generation.")
	}

	err = h.generateChapterOutlines(bookPath, bookState, &history, historyPath)
	if err != nil {
		return err
	}

	err = h.generateDrafts(bookPath, bookState, &history, historyPath)
	if err != nil {
		return err
	}

	h.Logger.Info(fmt.Sprintf("Book processing completed for topic: %s", topic))

	err = h.FileManager.SaveHistoryToFile(history, historyPath)
	if err != nil {
		return fmt.Errorf("failed to save conversation history: %w", err)
	}

	return nil
}

func (h *BookCommandHandler) generateBookOutline(topic, bookPath string, bookState *state.State, history *models.ConversationHistory, historyPath string) error {
	err := os.MkdirAll(bookPath, os.ModePerm)
	if err != nil {
		return h.handleError("failed to create book directory", err)
	}

	bookOutline := outline.NewOutline(topic)

	h.Logger.Info("Generating book outline...")

	prompt, err := h.WritingAgent.GenerateOutline(topic)
	if err != nil {
		return h.handleError("failed to generate outline prompt", err)
	}

	history.AddMessage("user", prompt)
	outlineContent, err := h.WritingAgent.SendMessage(history, h.FileManager, historyPath, false)
	if err != nil {
		if !h.ErrorHandler.HandleError(h.handleError("failed to generate book outline", err)) {
			return fmt.Errorf("retry attempts exhausted")
		}
	}

	h.Logger.Info(fmt.Sprintf("Generated outline content:\n%s", outlineContent))

	var yamlCheck map[string]interface{}
	err = yaml.Unmarshal([]byte(outlineContent), &yamlCheck)
	if err != nil {
		return h.handleError("failed to validate generated outline", err)
	}

	err = yaml.Unmarshal([]byte(outlineContent), bookOutline)
	if err != nil {
		return h.handleError("failed to parse generated outline", err)
	}

	err = bookOutline.Save(bookPath)
	if err != nil {
		return h.handleError("failed to save outline", err)
	}

	bookState.OutlineGenerated = true
	for _, chapter := range bookOutline.Chapters {
		chapterState := state.ChapterState{Title: chapter.Title, OutlineGenerated: true, DraftGenerated: false}
		for _, section := range chapter.Sections {
			sectionState := state.SectionState{Title: section.Title, OutlineGenerated: true, DraftGenerated: false}
			for _, subsection := range section.Subsections {
				sectionState.Subsections = append(sectionState.Subsections, state.SubsectionState{Title: subsection.Title})
			}
			chapterState.Sections = append(chapterState.Sections, sectionState)
		}
		bookState.Chapters = append(bookState.Chapters, chapterState)
	}

	err = bookState.Save(bookPath)
	if err != nil {
		return h.handleError("failed to save state", err)
	}

	h.Logger.Info(fmt.Sprintf("Book outline generated for topic: %s", topic))
	return nil
}

func (h *BookCommandHandler) generateChapterOutlines(bookPath string, bookState *state.State, history *models.ConversationHistory, historyPath string) error {
	for i, chapterState := range bookState.Chapters {
		if chapterState.OutlineGenerated && chapterState.DraftGenerated {
			continue
		}

		if chapterState.OutlineGenerated {
			h.Logger.Info(fmt.Sprintf("Outline for chapter %s is already generated, skipping.", chapterState.Title))
			continue
		}

		h.Logger.Info(fmt.Sprintf("Generating outline for chapter: %s", chapterState.Title))

		prompt, err := h.WritingAgent.GenerateChapterOutline(chapterState.Title)
		if err != nil {
			return h.handleError("failed to generate chapter outline prompt", err)
		}

		history.AddMessage("user", prompt)
		chapterOutlineContent, err := h.WritingAgent.SendMessage(history, h.FileManager, historyPath, false)
		if err != nil {
			if !h.ErrorHandler.HandleError(h.handleError("failed to generate chapter outline", err)) {
				return fmt.Errorf("retry attempts exhausted")
			}
		}

		h.Logger.Info(fmt.Sprintf("Generated chapter outline content:\n%s", chapterOutlineContent))

		chapterPath := filepath.Join(bookPath, fmt.Sprintf("ch%d", i+1))
		err = os.MkdirAll(chapterPath, os.ModePerm)
		if err != nil {
			return h.handleError("failed to create chapter directory", err)
		}

		chapterOutline := outline.NewChapterOutline(chapterState.Title)

		var yamlCheck map[string]interface{}
		err = yaml.Unmarshal([]byte(chapterOutlineContent), &yamlCheck)
		if err != nil {
			return h.handleError("failed to validate generated chapter outline", err)
		}

		err = yaml.Unmarshal([]byte(chapterOutlineContent), &chapterOutline)
		if err != nil {
			return h.handleError("failed to parse generated chapter outline", err)
		}

		err = chapterOutline.Save(chapterPath)
		if err != nil {
			return h.handleError("failed to save chapter outline", err)
		}

		if len(chapterState.Sections) < len(chapterOutline.Sections) {
			chapterState.Sections = make([]state.SectionState, len(chapterOutline.Sections))
		}

		chapterState.OutlineGenerated = true
		for j := range chapterOutline.Sections {
			chapterState.Sections[j].Title = chapterOutline.Sections[j].Title
			chapterState.Sections[j].OutlineGenerated = true

			if len(chapterState.Sections[j].Subsections) < len(chapterOutline.Sections[j].Subsections) {
				chapterState.Sections[j].Subsections = make([]state.SubsectionState, len(chapterOutline.Sections[j].Subsections))
			}

			for k := range chapterOutline.Sections[j].Subsections {
				chapterState.Sections[j].Subsections[k].Title = chapterOutline.Sections[j].Subsections[k].Title
			}
		}

		bookState.Chapters[i].OutlineGenerated = true
		err = bookState.Save(bookPath)
		if err != nil {
			return h.handleError("failed to save state", err)
		}

		history.AddMessage("assistant", chapterOutlineContent)
	}

	h.Logger.Info(fmt.Sprintf("All chapter outlines generated for book: %s", bookPath))
	return nil
}

func (h *BookCommandHandler) generateDrafts(bookPath string, bookState *state.State, history *models.ConversationHistory, historyPath string) error {
	for i, chapterState := range bookState.Chapters {
		if !chapterState.OutlineGenerated || chapterState.DraftGenerated {
			continue
		}

		chapterPath := filepath.Join(bookPath, fmt.Sprintf("ch%d", i+1))
		for j, section := range chapterState.Sections {
			if !section.OutlineGenerated || section.DraftGenerated {
				continue
			}

			h.Logger.Info(fmt.Sprintf("Generating draft for section: %s", section.Title))
			subsections := make([]outline.Subsection, len(section.Subsections))
			for k, subsection := range section.Subsections {
				subsections[k] = outline.Subsection{Title: subsection.Title}
			}
			sectionContent := outline.Section{
				Title:       section.Title,
				Subsections: subsections,
			}

			prompt, err := h.WritingAgent.GenerateSectionContent(sectionContent)
			if err != nil {
				return h.handleError("failed to generate section content prompt", err)
			}

			history.AddMessage("user", prompt)
			content, err := h.WritingAgent.SendMessage(history, h.FileManager, historyPath, true)
			if err != nil {
				if !h.ErrorHandler.HandleError(h.handleError("failed to generate section content", err)) {
					return fmt.Errorf("retry attempts exhausted")
				}
			}

			sectionPath := filepath.Join(chapterPath, fmt.Sprintf("section%d", j+1))
			err = os.MkdirAll(sectionPath, os.ModePerm)
			if err != nil {
				return h.handleError("failed to create section directory", err)
			}

			err = h.FileManager.SaveSectionContent(content, filepath.Join(sectionPath, "draft.md"))
			if err != nil {
				return h.handleError("failed to save section content", err)
			}

			chapterState.Sections[j].DraftGenerated = true
			err = bookState.Save(bookPath)
			if err != nil {
				return h.handleError("failed to save state", err)
			}

			// Add reference to saved content in the history
			history.AddMessage("assistant", fmt.Sprintf("Content saved to %s", filepath.Join(sectionPath, "draft.md")))
		}

		chapterState.DraftGenerated = true
		err := bookState.Save(bookPath)
		if err != nil {
			return h.handleError("failed to save state", err)
		}
	}

	h.Logger.Info(fmt.Sprintf("All drafts generated for book: %s", bookPath))
	return nil
}
