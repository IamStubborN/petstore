package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/IamStubborN/petstore/workers/api/mware"
	"gopkg.in/validator.v2"

	"github.com/IamStubborN/petstore/db"
	"github.com/IamStubborN/petstore/db/models"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

func StoreHandlers(r chi.Router) {
	r.Use(mware.JWT)
	r.Get("/inventory", inventoriesByStatus)
	r.Post("/order", createOrder)
	r.Get("/order/{order_ID}", findOrderByID)
	r.Delete("/order/{order_ID}", deletePurchaseByID)
}

func inventoriesByStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := genContext(r)
	defer cancel()

	storeDI := db.GetStoreDI()
	inventories, err := storeDI.GetInventories(ctx)
	if err != nil {
		respond(w, err, http.StatusInternalServerError, "providers error")
		return
	}

	data, err := json.Marshal(&inventories)
	if err != nil {
		respond(w, errors.Wrap(err, "can't marshal inventories to json"),
			http.StatusInternalServerError, "providers error")
		return
	}

	if _, err := w.Write(data); err != nil {
		respond(w, errors.Wrap(err, "can't write response"),
			http.StatusInternalServerError, "providers error")
		return
	}
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := genContext(r)
	defer cancel()

	if r.Body == nil {
		respond(w, errors.New("request body is nil"),
			http.StatusBadRequest, "invalid order")
		return
	}
	defer checkErrors(r.Body.Close)

	var order models.Order
	bytesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respond(w, errors.Wrap(err, "can't read request body"),
			http.StatusBadRequest, "invalid order")
		return
	}

	if err = order.UnmarshalJSON(bytesBody); err != nil {
		respond(w, errors.Wrap(err, "can't decode request body to pet"),
			http.StatusBadRequest, "invalid order")
		return
	}

	if err = validator.Validate(order); err != nil {
		respond(w, errors.Wrap(err, "can't validate order from body"),
			http.StatusMethodNotAllowed, "invalid input")
		return
	}

	orderQI := db.GetStoreDI()
	createdOrder, err := orderQI.CreateOrder(ctx, &order)
	if err != nil {
		respond(w, err, http.StatusBadRequest, "invalid order")
		return
	}

	data, err := createdOrder.MarshalJSON()
	if err != nil {
		respond(w, errors.Wrap(err, "can't marshal data"),
			http.StatusInternalServerError, "providers error")
		return
	}

	if _, err := w.Write(data); err != nil {
		respond(w, errors.Wrap(err, "can't write response"),
			http.StatusInternalServerError, "providers error")
		return
	}
}

func findOrderByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := genContext(r)
	defer cancel()

	slug := strings.TrimPrefix(r.URL.Path, "/api/v2/store/order/")
	id, err := strconv.ParseInt(slug, 10, 64)
	if err != nil {
		respond(w, errors.Wrapf(err, "can't cast slug to int [%s]", slug),
			http.StatusBadRequest, "invalid ID supplied")
		return
	}

	orderQI := db.GetStoreDI()
	order, err := orderQI.FindOrderByID(ctx, id)
	if err != nil {
		respond(w, err, http.StatusBadRequest, "order not found")
		return
	}

	data, err := order.MarshalJSON()
	if err != nil {
		respond(w, errors.Wrap(err, "can't marshal data"),
			http.StatusInternalServerError, "providers error")
		return
	}

	if _, err := w.Write(data); err != nil {
		respond(w, errors.Wrap(err, "can't write response"),
			http.StatusInternalServerError, "providers error")
		return
	}
}

func deletePurchaseByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := genContext(r)
	defer cancel()

	slug := strings.TrimPrefix(r.URL.Path, "/api/v2/store/order/")
	id, err := strconv.ParseInt(slug, 10, 64)
	if err != nil {
		respond(w, errors.Wrapf(err, "can't cast slug to int [%s]", slug),
			http.StatusBadRequest, "invalid ID supplied")
		return
	}

	orderQI := db.GetStoreDI()
	if err := orderQI.DeleteOrderByID(ctx, id); err != nil {
		respond(w, err, http.StatusNotFound, "order not found")
		return
	}

	respond(w, nil, http.StatusOK, "success")
}
