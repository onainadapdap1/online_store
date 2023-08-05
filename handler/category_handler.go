package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onainadapdap1/online_store/dtos"
	"github.com/onainadapdap1/online_store/models"
	"github.com/onainadapdap1/online_store/service"
	"github.com/onainadapdap1/online_store/utils"
)

type CategoryHandlerInterface interface {
	CreateCategory(c *gin.Context)
	UpdateCategory(c *gin.Context)
	FindBySlug(c *gin.Context)
	FindAllCategory(c *gin.Context)
	DeleteCategoryByID(c *gin.Context)
}

type categoryHandler struct {
	service service.CategoryServiceInterface
}

func NewCategoryHandler(service service.CategoryServiceInterface) CategoryHandlerInterface {
	return &categoryHandler{service: service}
}

func (h *categoryHandler) CreateCategory(c *gin.Context) {
	name := c.PostForm("name")
	description := c.PostForm("description")
	file, err := c.FormFile("image")
	if err != nil {
		response := utils.APIResponse("Failed to create category image", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	userId := currentUser.ID
	// filePath := fmt.Sprintf("static/images/categories/%d-%s", userId, file.Filename)

	fileName := fmt.Sprintf("%d-%s", userId, file.Filename)

	dirPath := filepath.Join(".", "static", "images", "categories")
	filePath := filepath.Join(dirPath, fileName)
	// Create directory if does not exist
	if _, err = os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			response := utils.APIResponse("Failed to upload category image", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusInternalServerError, response)
			return
		}
	}
	// Create file that will hold the image
	outputFile, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// Open the temporary file that contains the uploaded image
	inputFile, err := file.Open()
	if err != nil {
		response := utils.APIResponse("Failed to upload category image", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusOK, response)
	}
	defer inputFile.Close()

	// Copy the temporary image to the permanent location outputFile
	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		log.Fatal(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		return
	}

	createCategoryInput := dtos.CreateCategoryInput{
		User:        currentUser,
		Name:        name,
		Description: description,
		ImageURL:    filePath,
	}

	newCategory, err := h.service.CreateCategory(createCategoryInput)
	if err != nil {
		log.Printf("failed to create category: %v", err)
		response := utils.APIResponse("Failed to create category", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("Success to create category", http.StatusOK, "success", dtos.FormatCategory(newCategory))
	c.JSON(http.StatusOK, response)
}

func (h *categoryHandler) UpdateCategory(c *gin.Context) {
	var inputSlug dtos.GetCategoryDetailInput

	err := c.ShouldBindUri(&inputSlug)
	if err != nil {
		response := utils.APIResponse("Failed to get category slug", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData dtos.CreateCategoryInput

	name := c.PostForm("name")
	description := c.PostForm("description")

	category, err := h.service.FindBySlug(inputSlug)
	if err != nil {
		response := utils.APIResponse("Failed find by slug", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	userId := currentUser.ID

	// handle file after get data
	file, err := c.FormFile("image")
	if err != nil {
		// use existing image url if file is not found
		inputData.ImageURL = category.ImageURL
		response := utils.APIResponse("Failed to upload file", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	} else {
		// remove the old image file from the static folder,
		if category.ImageURL != "" {
			// oldFilename := strings.TrimPrefix(category.ImageURL, "/static/images/categories/")
			oldFilename := filepath.Base(category.ImageURL)
			if err := os.Remove("static/images/categories/" + oldFilename); err != nil {
				log.Printf("Failed to remove old filename: %v", err)
				response := utils.APIResponse(fmt.Sprintf("Failed to remove old filename: %v", err), http.StatusInternalServerError, "error", nil)
				c.JSON(http.StatusInternalServerError, response)
				return
			}
		}

		fileName := fmt.Sprintf("%d-%s", userId, file.Filename)

		dirPath := filepath.Join(".", "static", "images", "categories")
		filePath := filepath.Join(dirPath, fileName)
		// Create directory if does not exist
		if _, err = os.Stat(dirPath); os.IsNotExist(err) {
			err = os.MkdirAll(dirPath, 0755)
			if err != nil {
				response := utils.APIResponse("Failed to upload category image", http.StatusBadRequest, "error", nil)
				c.JSON(http.StatusInternalServerError, response)
				return
			}
		}
		// Create file that will hold the image
		outputFile, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer outputFile.Close()

		// Open the temporary file that contains the uploaded image
		inputFile, err := file.Open()
		if err != nil {
			response := utils.APIResponse("Failed to upload category image", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusOK, response)
		}
		defer inputFile.Close()

		// Copy the temporary image to the permanent location outputFile
		_, err = io.Copy(outputFile, inputFile)
		if err != nil {
			log.Fatal(err)
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		inputData.ImageURL = filePath
	}
	inputData.User = currentUser
	inputData.Name = name
	inputData.Description = description

	updatedCategory, err := h.service.UpdateCategory(inputSlug, inputData)
	if err != nil {
		log.Printf("failed to update category: %v", err)
		response := utils.APIResponse(fmt.Sprintf("failed to update category: %v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("sucess to update category", http.StatusOK, "success", dtos.FormatCategory(updatedCategory))
	c.JSON(http.StatusOK, response)
}

func (h *categoryHandler) FindBySlug(c *gin.Context) {
	var input dtos.GetCategoryDetailInput
	// var category models.Category
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := utils.APIResponse("failed to get detail input", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	category, err := h.service.FindBySlug(input)
	if err != nil {
		response := utils.APIResponse("failed to get detail category", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("success get category detail", http.StatusOK, "success", dtos.FormateCategoryDetail(category))
	c.JSON(http.StatusOK, response)
}

func (h *categoryHandler) FindAllCategory(c *gin.Context) {
	categories, err := h.service.FindAllCategory()
	if err != nil {
		response := utils.APIResponse("failed to get all categories", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("list of categories", http.StatusOK, "success", dtos.FormatCategories(categories))
	c.JSON(http.StatusOK, response)
}


func (h *categoryHandler) DeleteCategoryByID(c *gin.Context) {
	// var input models.GetCategoryDetailInput

	param := c.Param("id")
	categoryID, err := strconv.Atoi(param)
	if err != nil {
		response := utils.APIResponse("failed to get detail input", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	category, err := h.service.FindByCategoryID(uint(categoryID))
	if err != nil {
		response := utils.APIResponse("failed to get detail category", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if category.ImageURL != "" {
		oldFilename := filepath.Base(category.ImageURL)
		if err := os.Remove("static/images/categories/" + oldFilename); err != nil {
			log.Printf("Failed to remove old filename: %v", err)
			response := utils.APIResponse(fmt.Sprintf("Failed to remove old filename: %v", err), http.StatusInternalServerError, "error", nil)
			c.JSON(http.StatusInternalServerError, response)
			return
		}
	}

	err = h.service.DeleteCategory(category)
	if err != nil {
		response := utils.APIResponse("failed to delete category", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := utils.APIResponse("Success to delete category", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
