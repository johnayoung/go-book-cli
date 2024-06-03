package handlers

import (
	"fmt"
	"go-book-ai/internal/agents"
	"go-book-ai/internal/errors"
	"go-book-ai/internal/file"
	"go-book-ai/internal/logger"
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

func (h *BookCommandHandler) ProcessBook(topic string) error {
	folderName := utils.CleanName(topic)
	bookPath := filepath.Join("books", folderName)
	h.Logger.Info(fmt.Sprintf("Processing book folder: %s", bookPath))

	stateFilePath := filepath.Join(bookPath, "state.yaml")

	// Load state
	bookState, err := h.FileManager.LoadState(stateFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			bookState = state.NewState()
		} else {
			return fmt.Errorf("failed to load state: %v", err)
		}
	}

	h.Logger.Info(fmt.Sprintf("Loaded state: %+v", bookState))

	// Check if the outline has already been generated
	if !bookState.OutlineGenerated {
		err := h.generateBookOutline(topic, bookPath, bookState)
		if err != nil {
			return err
		}
		h.Logger.Info(fmt.Sprintf("State after generating book outline: %+v", bookState))
		err = h.FileManager.SaveState(stateFilePath, bookState)
		if err != nil {
			return fmt.Errorf("failed to save state after outline generation: %w", err)
		}
	} else {
		h.Logger.Info("Outline already generated, skipping outline generation.")
	}

	err = h.generateChapterOutlines(bookPath, bookState)
	if err != nil {
		return err
	}
	h.Logger.Info(fmt.Sprintf("State after generating chapter outlines: %+v", bookState))
	err = h.FileManager.SaveState(stateFilePath, bookState)
	if err != nil {
		return fmt.Errorf("failed to save state after chapter outlines generation: %w", err)
	}

	err = h.generateDrafts(bookPath, bookState)
	if err != nil {
		return err
	}
	h.Logger.Info(fmt.Sprintf("State after generating drafts: %+v", bookState))
	err = h.FileManager.SaveState(stateFilePath, bookState)
	if err != nil {
		return fmt.Errorf("failed to save state after drafts generation: %w", err)
	}

	h.Logger.Info(fmt.Sprintf("Book processing completed for topic: %s", topic))
	return nil
}

func (h *BookCommandHandler) generateBookOutline(topic, bookPath string, bookState *state.State) error {
	err := os.MkdirAll(bookPath, os.ModePerm)
	if err != nil {
		return h.handleError("failed to create book directory", err)
	}

	h.Logger.Info("Generating book outline...")

	prompt, err := h.WritingAgent.GenerateOutline(topic)
	if err != nil {
		return h.handleError("failed to generate outline prompt", err)
	}

	bookState.MessageHistory = append(bookState.MessageHistory, state.Message{Role: "user", Content: prompt})
	outlineContent, err := h.WritingAgent.SendMessage(&bookState.MessageHistory, h.FileManager, false)
	if err != nil {
		if !h.ErrorHandler.HandleError(h.handleError("failed to generate book outline", err)) {
			return fmt.Errorf("retry attempts exhausted")
		}
	}

	h.Logger.Info(fmt.Sprintf("Generated outline content:\n%s", outlineContent))

	var outline struct {
		Title    string               `yaml:"title"`
		Chapters []state.ChapterState `yaml:"chapters"`
	}

	err = yaml.Unmarshal([]byte(outlineContent), &outline)
	if err != nil {
		return h.handleError("failed to parse generated outline", err)
	}

	bookState.OutlineGenerated = true
	bookState.Chapters = outline.Chapters

	for i := range bookState.Chapters {
		bookState.Chapters[i].OutlineGenerated = false
		bookState.Chapters[i].DraftGenerated = false
	}

	bookState.MessageHistory = append(bookState.MessageHistory, state.Message{Role: "assistant", Content: outlineContent})

	h.Logger.Info(fmt.Sprintf("Book outline generated for topic: %s", topic))
	return nil
}

func (h *BookCommandHandler) generateChapterOutlines(bookPath string, bookState *state.State) error {
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

		bookState.MessageHistory = append(bookState.MessageHistory, state.Message{Role: "user", Content: prompt})
		chapterOutlineContent, err := h.WritingAgent.SendMessage(&bookState.MessageHistory, h.FileManager, false)
		if err != nil {
			if !h.ErrorHandler.HandleError(h.handleError("failed to generate chapter outline", err)) {
				return fmt.Errorf("retry attempts exhausted")
			}
		}

		h.Logger.Info(fmt.Sprintf("Generated chapter outline content:\n%s", chapterOutlineContent))

		var chapterOutline state.ChapterState

		err = yaml.Unmarshal([]byte(chapterOutlineContent), &chapterOutline)
		if err != nil {
			return h.handleError("failed to parse generated chapter outline", err)
		}

		bookState.Chapters[i].OutlineGenerated = true
		bookState.Chapters[i].Sections = chapterOutline.Sections

		for j := range bookState.Chapters[i].Sections {
			bookState.Chapters[i].Sections[j].OutlineGenerated = false
			bookState.Chapters[i].Sections[j].DraftGenerated = false
		}

		bookState.MessageHistory = append(bookState.MessageHistory, state.Message{Role: "assistant", Content: chapterOutlineContent})
	}

	h.Logger.Info(fmt.Sprintf("All chapter outlines generated for book: %s", bookPath))
	return nil
}

func (h *BookCommandHandler) generateDrafts(bookPath string, bookState *state.State) error {
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

			bookState.MessageHistory = append(bookState.MessageHistory, state.Message{Role: "user", Content: prompt})
			content, err := h.WritingAgent.SendMessage(&bookState.MessageHistory, h.FileManager, true)
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
			err = h.FileManager.SaveState(filepath.Join(bookPath, "state.yaml"), bookState)
			if err != nil {
				return h.handleError("failed to save state", err)
			}

			// Add reference to saved content in the history
			bookState.MessageHistory = append(bookState.MessageHistory, state.Message{Role: "assistant", Content: fmt.Sprintf("Content saved to %s", filepath.Join(sectionPath, "draft.md"))})
		}

		chapterState.DraftGenerated = true
		err := h.FileManager.SaveState(filepath.Join(bookPath, "state.yaml"), bookState)
		if err != nil {
			return h.handleError("failed to save state", err)
		}
	}

	h.Logger.Info(fmt.Sprintf("All drafts generated for book: %s", bookPath))
	return nil
}

func (h *BookCommandHandler) handleError(message string, err error) error {
	h.ErrorHandler.LogError(fmt.Errorf("%s: %w", message, err))
	return fmt.Errorf("%s: %w", message, err)
}
