package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/IamStubborN/petstore/db"
	"github.com/IamStubborN/petstore/db/models"
	"github.com/IamStubborN/petstore/workers/api/mware"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"gopkg.in/validator.v2"
)

func PetHandlers(r chi.Router) {
	r.Use(mware.JWT)
	r.Post("/", addPetToStore)
	r.Put("/", updatePetInStore)
	r.Get("/findByStatus", findPetsByStatus)
	r.Get("/{petID}", getPetByID)
	r.Post("/{petID}", updatePetByID)
	r.Delete("/{petID}", deletePetByID)
	r.Post("/{petID}/uploadImage", updatePetPhotosByID)
}

func addPetToStore(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := genContext(r)
	defer cancel()

	if r.Body == nil {
		respond(w, errors.New("request body is nil"),
			http.StatusMethodNotAllowed, "invalid input")
		return
	}
	defer checkErrors(r.Body.Close)

	var pet models.Pet
	bytesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respond(w, errors.Wrap(err, "can't read request body"),
			http.StatusMethodNotAllowed, "invalid input")
		return
	}

	if err = pet.UnmarshalJSON(bytesBody); err != nil {
		respond(w, errors.Wrap(err, "can't decode request body to pet"),
			http.StatusMethodNotAllowed, "invalid input")
		return
	}

	if err = validator.Validate(pet); err != nil {
		respond(w, errors.Wrap(err, "can't validate pet from body"),
			http.StatusMethodNotAllowed, "invalid input")
		return
	}

	petDI := db.GetPetDI()
	addedPet, err := petDI.AddPetToStore(ctx, &pet)
	if err != nil {
		respond(w, err, http.StatusMethodNotAllowed, "invalid input")
		return
	}

	data, err := addedPet.MarshalJSON()
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

func updatePetInStore(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := genContext(r)
	defer cancel()

	if r.Body == nil {
		respond(w, errors.New("request body is nil"),
			http.StatusMethodNotAllowed, "validation exception")
		return
	}
	defer checkErrors(r.Body.Close)

	var pet models.Pet
	bytesBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respond(w, errors.Wrap(err, "can't read request body"),
			http.StatusMethodNotAllowed, "validation exception")
		return
	}

	if err = pet.UnmarshalJSON(bytesBody); err != nil {
		respond(w, errors.Wrap(err, "can't decode request body to pet"),
			http.StatusBadRequest, "invalid ID supplied")
		return
	}

	if err = validator.Validate(pet); err != nil {
		respond(w, errors.Wrap(err, "can't validate pet from body"),
			http.StatusMethodNotAllowed, "invalid input")
		return
	}

	petDI := db.GetPetDI()
	updatedPet, err := petDI.UpdatePetInStoreByBody(ctx, &pet)
	if err != nil {
		respond(w, errors.Wrap(err, "can't update pet in store"),
			http.StatusNotFound, "pet not found")
		return
	}

	data, err := updatedPet.MarshalJSON()
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

func findPetsByStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := genContext(r)
	defer cancel()

	statusQuery := r.URL.Query().Get("status")
	if statusQuery == "" {
		respond(w, errors.New("no status in query"),
			http.StatusBadRequest, "invalid status value")
		return
	}

	status := strings.Split(statusQuery, ",")
	PetDI := db.GetPetDI()
	pets, err := PetDI.FindPetsByStatus(ctx, status)
	if err != nil {
		respond(w, err, http.StatusBadRequest, "invalid status value")
		return
	}

	data, err := pets.MarshalJSON()
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

func getPetByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := genContext(r)
	defer cancel()

	slug := strings.TrimPrefix(r.URL.Path, "/api/v2/pet/")
	id, err := strconv.ParseInt(slug, 10, 64)
	if err != nil {
		respond(w, errors.Wrapf(err, "can't cast slug to int [%s]", slug),
			http.StatusBadRequest, "invalid ID supplied")
		return
	}

	PetDI := db.GetPetDI()
	pet, err := PetDI.GetPetByID(ctx, id)
	if err != nil {
		respond(w, err, http.StatusNotFound, "pet not found")
		return
	}

	data, err := pet.MarshalJSON()
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

func updatePetByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := genContext(r)
	defer cancel()

	if err := r.ParseForm(); err != nil {
		respond(w, errors.New("can't parse form"),
			http.StatusMethodNotAllowed, "invalid input")
		return
	}

	slug := strings.TrimPrefix(r.URL.Path, "/api/v2/pet/")
	id, err := strconv.ParseInt(slug, 10, 64)
	if err != nil {
		respond(w, errors.Wrapf(err, "can't cast slug to int [%s]", slug),
			http.StatusBadRequest, "invalid input")
		return
	}

	name := r.FormValue("name")
	status := r.FormValue("status")

	if isNotValid(name) && isNotValid(status) {
		respond(w, fmt.Errorf("bad input %s or %s", name, status),
			http.StatusMethodNotAllowed, "invalid input")
		return
	}

	PetDI := db.GetPetDI()
	if err := PetDI.UpdatePetInStoreByForm(ctx, id, name, status); err != nil {
		respond(w, errors.Wrap(err, "can't update pet"),
			http.StatusMethodNotAllowed, "invalid input")
		return
	}

	respond(w, nil, http.StatusOK, "success")
}

func deletePetByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := genContext(r)
	defer cancel()

	if err := r.ParseForm(); err != nil {
		respond(w, errors.New("can't parse form"),
			http.StatusMethodNotAllowed, "invalid input")
		return
	}

	slug := strings.TrimPrefix(r.URL.Path, "/api/v2/pet/")
	id, err := strconv.ParseInt(slug, 10, 64)
	if err != nil {
		respond(w, errors.Wrapf(err, "can't cast slug to int [%s]", slug),
			http.StatusBadRequest, "invalid ID supplied")
		return
	}

	PetDI := db.GetPetDI()
	if err := PetDI.DeletePetByID(ctx, id); err != nil {
		respond(w, err, http.StatusBadRequest, "invalid ID supplied")
		return
	}

	respond(w, nil, http.StatusOK, "success")
}

func updatePetPhotosByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := genContext(r)
	defer cancel()

	if err := r.ParseForm(); err != nil {
		respond(w, errors.New("can't parse form"),
			http.StatusMethodNotAllowed, "invalid input")
		return
	}

	slug := strings.TrimPrefix(r.URL.Path, "/api/v2/pet/")
	slug = strings.TrimSuffix(slug, "/uploadImage")
	id, err := strconv.ParseInt(slug, 10, 64)
	if err != nil {
		respond(w, errors.Wrapf(err, "can't cast slug to int [%s]", slug),
			http.StatusBadRequest, "invalid input")
		return
	}

	images := strings.Split(r.FormValue("images"), ",")

	PetDI := db.GetPetDI()
	if err := PetDI.UpdatePetPhotosByID(ctx, id, images); err != nil {
		respond(w, err, http.StatusBadRequest, "invalid input")
		return
	}

	respond(w, nil, http.StatusOK, "success")
}
