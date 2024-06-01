package handlers

import (
	"fmt"
	"os"
	"path/filepath"

	"go-book-ai/internal/agents"
	"go-book-ai/internal/errors"
	"go-book-ai/internal/file"
	"go-book-ai/internal/logger"
	"go-book-ai/internal/outline"
	"go-book-ai/internal/state"
	"go-book-ai/internal/utils"

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

	bookState, err := state.LoadState(bookPath)
	if err != nil {
		if os.IsNotExist(err) {
			bookState = state.NewState()
		} else {
			return fmt.Errorf("failed to load state: %v", err)
		}
	}

	if !bookState.OutlineGenerated {
		err := h.generateBookOutline(topic, bookPath, bookState)
		if err != nil {
			return err
		}
	}

	err = h.generateChapterOutlines(bookPath, bookState)
	if err != nil {
		return err
	}

	h.Logger.Info(fmt.Sprintf("Book processing completed for topic: %s", topic))
	return nil
}

func (h *BookCommandHandler) generateBookOutline(topic, bookPath string, bookState *state.State) error {
	err := os.MkdirAll(bookPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create book directory: %v", err)
	}

	bookOutline := outline.NewOutline(topic)

	h.Logger.Info("Generating book outline...")
	outlineContent, err := h.WritingAgent.GenerateOutline(topic)
	if err != nil {
		h.ErrorHandler.LogError(fmt.Errorf("failed to generate book outline: %v", err))
		if !h.ErrorHandler.HandleError(err) {
			return fmt.Errorf("retry attempts exhausted")
		}
	}

	// Log the generated outline content to understand its format
	h.Logger.Info(fmt.Sprintf("Generated outline content:\n%s", outlineContent))

	// Validate the YAML format
	var yamlCheck map[string]interface{}
	err = yaml.Unmarshal([]byte(outlineContent), &yamlCheck)
	if err != nil {
		h.ErrorHandler.LogError(fmt.Errorf("failed to validate generated outline: %v", err))
		return fmt.Errorf("failed to validate generated outline: %v", err)
	}

	// Assuming the generated outline content is in YAML format,
	// we need to parse it and populate the `Outline` struct.
	err = yaml.Unmarshal([]byte(outlineContent), bookOutline)
	if err != nil {
		h.ErrorHandler.LogError(fmt.Errorf("failed to parse generated outline: %v", err))
		return fmt.Errorf("failed to parse generated outline: %v", err)
	}

	err = bookOutline.Save(bookPath)
	if err != nil {
		h.ErrorHandler.LogError(fmt.Errorf("failed to save outline: %v", err))
		return fmt.Errorf("failed to save outline: %v", err)
	}

	bookState.OutlineGenerated = true
	for _, chapter := range bookOutline.Chapters {
		chapterState := state.ChapterState{Title: chapter.Title, Generated: false}
		for _, section := range chapter.Sections {
			sectionState := state.SectionState{Title: section.Title, Generated: false}
			for _, subsection := range section.Subsections {
				sectionState.Subsections = append(sectionState.Subsections, state.SubsectionState{Title: subsection.Title, Generated: false})
			}
			chapterState.Sections = append(chapterState.Sections, sectionState)
		}
		bookState.Chapters = append(bookState.Chapters, chapterState)
	}
	err = bookState.Save(bookPath)
	if err != nil {
		h.ErrorHandler.LogError(fmt.Errorf("failed to save state: %v", err))
		return fmt.Errorf("failed to save state: %v", err)
	}

	h.Logger.Info(fmt.Sprintf("Book outline generated for topic: %s", topic))
	return nil
}

func (h *BookCommandHandler) generateChapterOutlines(bookPath string, bookState *state.State) error {
	for i, chapterState := range bookState.Chapters {
		if chapterState.Generated {
			continue
		}

		h.Logger.Info(fmt.Sprintf("Generating outline for chapter: %s", chapterState.Title))
		chapterOutlineContent, err := h.WritingAgent.GenerateChapterOutline(chapterState.Title)
		if err != nil {
			h.ErrorHandler.LogError(fmt.Errorf("failed to generate chapter outline: %v", err))
			if !h.ErrorHandler.HandleError(err) {
				return fmt.Errorf("retry attempts exhausted")
			}
		}

		h.Logger.Info(fmt.Sprintf("Generated chapter outline content:\n%s", chapterOutlineContent))

		chapterPath := filepath.Join(bookPath, fmt.Sprintf("ch%d", i+1))
		err = os.MkdirAll(chapterPath, os.ModePerm)
		if err != nil {
			h.ErrorHandler.LogError(fmt.Errorf("failed to create chapter directory: %v", err))
			return fmt.Errorf("failed to create chapter directory: %v", err)
		}

		chapterOutline := outline.NewChapterOutline(chapterState.Title)

		// Validate the YAML format
		var yamlCheck map[string]interface{}
		err = yaml.Unmarshal([]byte(chapterOutlineContent), &yamlCheck)
		if err != nil {
			h.ErrorHandler.LogError(fmt.Errorf("failed to validate generated chapter outline: %v", err))
			return fmt.Errorf("failed to validate generated chapter outline: %v", err)
		}

		err = yaml.Unmarshal([]byte(chapterOutlineContent), &chapterOutline)
		if err != nil {
			h.ErrorHandler.LogError(fmt.Errorf("failed to parse generated chapter outline: %v", err))
			return fmt.Errorf("failed to parse generated chapter outline: %v", err)
		}

		err = chapterOutline.Save(chapterPath)
		if err != nil {
			h.ErrorHandler.LogError(fmt.Errorf("failed to save chapter outline: %v", err))
			return fmt.Errorf("failed to save chapter outline: %v", err)
		}

		chapterState.Generated = true
		for j, section := range chapterOutline.Sections {
			chapterState.Sections[j].Title = section.Title
			chapterState.Sections[j].Generated = true
			for k, subsection := range section.Subsections {
				chapterState.Sections[j].Subsections[k].Title = subsection.Title
				chapterState.Sections[j].Subsections[k].Generated = true
			}
		}

		bookState.Chapters[i].Generated = true
		err = bookState.Save(bookPath)
		if err != nil {
			h.ErrorHandler.LogError(fmt.Errorf("failed to save state: %v", err))
			return fmt.Errorf("failed to save state: %v", err)
		}
	}

	h.Logger.Info(fmt.Sprintf("All chapter outlines generated for book: %s", bookPath))
	return nil
}
