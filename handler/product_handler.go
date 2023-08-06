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

type ProductHandlerInterface interface {
	CreateProduct(c *gin.Context)
	UpdateProduct(c *gin.Context)
	FindProductBySlug(c *gin.Context)
	FindAllProduct(c *gin.Context)
}

type productHandler struct {
	service service.ProductServiceInterface
}

func NewProductHandler(service service.ProductServiceInterface) ProductHandlerInterface {
	return &productHandler{service: service}
}

// Create Product godoc
// @Summary Create Product
// @Description Create a new Product with a given name, description and image file
// @Tags products
// @Accept mpfd
// @Produce json
// @Param name formData string true "Name of the Product"
// @Param description formData string true "Description of the product"
// @Param price formData number true "price of the product"
// @Param quantity formData int true "quantity of the product"
// @Param category_id formData int true "category id of the product"
// @Param image formData file true "Image file of the product"
// @Success 200 {object} dtos.ProductFormatter
// @Failure 400 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/products/product [post]
func (h *productHandler) CreateProduct(c *gin.Context) {
	// get data from request
	name := c.PostForm("name")
	description := c.PostForm("description")
	price := c.PostForm("price")
	priceInt, err := strconv.Atoi(price)
	if err != nil {
		response := utils.APIResponse("failed to convert price", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	quantity := c.PostForm("quantity")
	quantityInt, err := strconv.Atoi(quantity)
	if err != nil {
		response := utils.APIResponse("failed to convert price", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	category_id := c.PostForm("category_id")
	categoryIdInt, err := strconv.Atoi(category_id)
	if err != nil {
		response := utils.APIResponse("failed to convert category id", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	file, err := c.FormFile("image")
	if err != nil {
		response := utils.APIResponse("failed to load image file", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	userId := currentUser.ID

	fileName := fmt.Sprintf("%d-%s", userId, file.Filename)

	dirPath := filepath.Join(".", "static", "images", "products")
	filePath := filepath.Join(dirPath, fileName)
	// create directory if doesn't exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			response := utils.APIResponse("failed to upload product image", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusBadRequest, response)
			return
		}
	}
	// create file that will hold the image
	outputFile, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// open the temporary file that contains the uploaded image
	inputFile, err := file.Open()
	if err != nil {
		response := utils.APIResponse("failed to open product input image", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	defer inputFile.Close()

	// copy the temporary image to the permanent location outputFile
	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		response := utils.APIResponse("failed to copy product input image to permanent location", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	category, err := h.service.GetCategoryByID(uint(categoryIdInt))
	if err != nil {
		response := utils.APIResponse("failed to get category by id", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	createProductInput := dtos.CreateProductInput{
		Name:        name,
		Description: description,
		Price:       float64(priceInt),
		Quantity:    quantityInt,
		CategoryID:  category.ID,
		ImageURL:    filePath,
		User:        currentUser,
		Category:    category,
	}

	newProduct, err := h.service.CreateProduct(createProductInput)
	if err != nil {
		log.Printf("failed to create product: %v", err)
		response := utils.APIResponse("failed to create produt", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("success to create product", http.StatusOK, "success", dtos.FormatProduct(newProduct))
	c.JSON(http.StatusOK, response)
}

// Update product godoc
// @Summary Update product
// @Description Update product
// @Tags products
// @Accept json,mpfd
// @Produce json
// @Param slug path string true "update category by slug"
// @Param name formData string true "Name of the Product"
// @Param description formData string true "Description of the product"
// @Param price formData number true "price of the product"
// @Param quantity formData int true "quantity of the product"
// @Param category_id formData int true "category id of the product"
// @Param image formData file true "Image file of the product"
// @Success 200 {object} dtos.ProductFormatter
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Security BearerAuth
// @Router /api/v1/products/product/{slug} [put]
func (h *productHandler) UpdateProduct(c *gin.Context) {
	var inputSlug dtos.GetProductDetailInput

	if err := c.ShouldBindUri(&inputSlug); err != nil {
		response := utils.APIResponse("faile to get product slug", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	name := c.PostForm("name")
	description := c.PostForm("description")
	price := c.PostForm("price")
	priceInt, err := strconv.Atoi(price)
	if err != nil {
		response := utils.APIResponse("failed to convert price", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	quantity := c.PostForm("quantity")
	quantityInt, err := strconv.Atoi(quantity)
	if err != nil {
		response := utils.APIResponse("failed to convert price", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	category_id := c.PostForm("category_id")
	categoryIdInt, err := strconv.Atoi(category_id)
	if err != nil {
		response := utils.APIResponse("failed to convert category id", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// file, err := c.FormFile("image")
	// if err != nil {
	// 	response := utils.APIResponse("failed to load image file", http.StatusBadRequest, "error", nil)
	// 	c.JSON(http.StatusBadRequest, response)
	// 	return
	// }
	category, err := h.service.GetCategoryByID(uint(categoryIdInt))
	if err != nil {
		response := utils.APIResponse("failed to get category by id", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// get product by slug
	product, err := h.service.FindProductBySlug(inputSlug)
	if err != nil {
		response := utils.APIResponse("Failed find product by slug", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)
	userId := currentUser.ID

	var inputData dtos.CreateProductInput

	file, err := c.FormFile("image")
	if err != nil {
		inputData.ImageURL = product.ImageURL
		response := utils.APIResponse("Failed to upload file", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	} else {
		// remove the old image file from the static folder
		if product.ImageURL != "" {
			oldFilename := filepath.Base(product.ImageURL)
			if err := os.Remove("static/images/products/" + oldFilename); err != nil {
				log.Printf("failed to remove old filename: %v", err)
				response := utils.APIResponse(fmt.Sprintf("Failed to remove old filename: %v", err), http.StatusInternalServerError, "error", nil)
				c.JSON(http.StatusInternalServerError, response)
				return
			}
		}
		fileName := fmt.Sprintf("%d-%s", userId, file.Filename)
		dirPath := filepath.Join(".", "static", "images", "products")
		filePath := filepath.Join(dirPath, fileName)

		// create file that will hold the image
		outputFile, err := os.Create(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer outputFile.Close()
		inputFile, err := file.Open()
		if err != nil {
			response := utils.APIResponse("Failed to upload product image", http.StatusBadRequest, "error", nil)
			c.JSON(http.StatusOK, response)
		}
		defer inputFile.Close()
		_, err = io.Copy(outputFile, inputFile)
		if err != nil {
			log.Fatal(err)
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		inputData.ImageURL = filePath
	}
	// data := models.CreateProductInput{}
	inputData.User = currentUser
	inputData.Category = category
	inputData.Name = name
	inputData.Description = description
	inputData.Price = float64(priceInt)
	inputData.Quantity = quantityInt
	inputData.CategoryID = uint(categoryIdInt)

	updatedProduct, err := h.service.UpdateProduct(inputSlug, inputData)
	if err != nil {
		log.Printf("failed to update product: %v", err)
		response := utils.APIResponse(fmt.Sprintf("failed to update product: %v", err), http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("success to update product", http.StatusOK, "success", dtos.FormatProduct(updatedProduct))
	c.JSON(http.StatusOK, response)
}


// Get product by slug godoc
// @Summary Get one product by slug
// @Description Get one product by slug
// @Tags products
// @Produce json
// @Param slug path string true "get product by slug"
// @Success 200 {object} dtos.ProductDetailFormatter{}
// @Failure 400 {object} utils.Response
// @Router /api/v1/products/product/{slug} [get]
func (h *productHandler) FindProductBySlug(c *gin.Context) {
	var input dtos.GetProductDetailInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := utils.APIResponse("failed to get detail input", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	product, err := h.service.FindProductBySlug(input)
	if err != nil {
		response := utils.APIResponse("failed to get detail product", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("success to get product detail", http.StatusOK, "success", dtos.FormatProductDetail(product))
	c.JSON(http.StatusOK, response)
}

// Get All products godoc
// @Summary Get all products
// @Description Get all products
// @Tags products
// @Produce json
// @Success 200 {object} []dtos.ProductDetailFormatter{}
// @Failure 400 {object} utils.Response
// @Router /api/v1/products [get]
func (h *productHandler) FindAllProduct(c *gin.Context) {
	products, err := h.service.FindAllProduct()
	if err != nil {
		response := utils.APIResponse("failed to get all products", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := utils.APIResponse("list of products", http.StatusOK, "success", dtos.FormatProducts(products))
	c.JSON(http.StatusOK, response)
}
