package api

import (
	DB "carbide-images-api/pkg/database"
	"carbide-images-api/pkg/objects"
	"database/sql"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func productNameFromPath(r *http.Request) string {
	productName := r.PathValue("name")
	return productName
}

// Responds with a JSON array of all products in the database
//
// Success Code: 200 OK
func getAllProductsHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		products, err := DB.GetAllProducts(db)
		if err != nil {
			httpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		err = sendAsJSON(w, products)
		if err != nil {
			log.Error(err)
		}
		return
	}
	return http.HandlerFunc(fn)
}

// Accepts a JSON payload of a new product and responds with the new JSON object after it's been successfully created in the database
//
// Success Code: 201 OK
func createProductHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var createdProduct objects.Product
		err := decodeJSONObject(w, r, &createdProduct)
		if err != nil {
			log.Error(err)
			return
		}
		err = DB.AddProduct(db, createdProduct)
		if err != nil {
			httpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		createdProduct, err = DB.GetProduct(db, *createdProduct.Name)
		if err != nil {
			httpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		log.WithFields(log.Fields{
			"product": *createdProduct.Name,
		}).Info("product has been successfully created")
		w.WriteHeader(http.StatusCreated)
		err = sendAsJSON(w, createdProduct)
		if err != nil {
			log.Error(err)
		}
		return
	}
	return http.HandlerFunc(fn)
}

// Responds with the JSON representation of a product
//
// Success Code: 200 OK
func getProductHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var retrievedProduct objects.Product
		productName := productNameFromPath(r)
		retrievedProduct, err := DB.GetProduct(db, productName)
		if err != nil {
			httpJSONError(w, err.Error(), http.StatusBadRequest)
			log.Error(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		err = sendAsJSON(w, retrievedProduct)
		if err != nil {
			log.Error(err)
		}
		return
	}
	return http.HandlerFunc(fn)
}

// Responds with the JSON representation of a product
//
// Success Code: 200 OK
func updateProductHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var updatedProduct objects.Product
		productName := productNameFromPath(r)
		err := decodeJSONObject(w, r, &updatedProduct)
		if err != nil {
			log.Error(err)
			return
		}
		err = DB.UpdateProduct(db, *updatedProduct.Name, productName)
		if err != nil {
			log.Error(err)
			return
		}
		updatedProduct, err = DB.GetProduct(db, *updatedProduct.Name)
		if err != nil {
			log.Error(err)
			return
		}
		log.WithFields(log.Fields{
			"product": *updatedProduct.Name,
		}).Info("product has been successfully updated")
		w.WriteHeader(http.StatusOK)
		err = sendAsJSON(w, updatedProduct)
		if err != nil {
			log.Error(err)
		}
		return
	}
	return http.HandlerFunc(fn)
}

// Deletes the product and responds with an empty payload
//
// Success Code: 204 No Content
func deleteProductHandler(db *sql.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		productName := productNameFromPath(r)
		err := DB.DeleteProduct(db, productName)
		if err != nil {
			httpJSONError(w, err.Error(), http.StatusInternalServerError)
			log.Error(err)
			return
		}
		log.WithFields(log.Fields{
			"product": productName,
		}).Info("product has been successfully deleted")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	return http.HandlerFunc(fn)
}
