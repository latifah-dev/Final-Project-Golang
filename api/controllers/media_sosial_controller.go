package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tifarin/fullstack/api/auth"
	"github.com/tifarin/fullstack/api/models"
	"github.com/tifarin/fullstack/api/responses"
	"github.com/tifarin/fullstack/api/utils/formaterror"
)

func (server *Server) CreateMediaSosial(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	mediasosial := models.Media_Sosial{}
	err = json.Unmarshal(body, &mediasosial)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	mediasosial.Prepare()
	err = mediasosial.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != mediasosial.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	mediasosialCreated, err := mediasosial.SaveMediaSosial(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, mediasosialCreated.ID))
	responses.JSON(w, http.StatusCreated, mediasosialCreated)
}

func (server *Server) GetMediaSosials(w http.ResponseWriter, r *http.Request) {

	mediasosial := models.Photo{}

	mediasosials, err := mediasosial.FindAllPhotos(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, mediasosials)
}

func (server *Server) GetMediaSosial(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	mediasosial := models.Media_Sosial{}

	mediasosialReceived, err := mediasosial.FindMediaSosialByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, mediasosialReceived)
}

func (server *Server) UpdateAMediaSosial(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the post id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the post exist
	mediasosial := models.Media_Sosial{}
	err = server.DB.Debug().Model(models.Media_Sosial{}).Where("id = ?", pid).Take(&mediasosial).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Comment not found"))
		return
	}

	// If a user attempt to update a post not belonging to him
	if uid != mediasosial.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data posted
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	mediasosialUpdate := models.Media_Sosial{}
	err = json.Unmarshal(body, &mediasosialUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != mediasosialUpdate.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	mediasosialUpdate.Prepare()
	err = mediasosialUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	mediasosialUpdate.ID = mediasosial.ID //this is important to tell the model the post id to update, the other update field are set above

	mediasosialUpdated, err := mediasosialUpdate.UpdateAMediaSosial(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, mediasosialUpdated)
}

func (server *Server) DeleteAMediaSosial(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid post id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the post exist
	mediasosial := models.Media_Sosial{}
	err = server.DB.Debug().Model(models.Media_Sosial{}).Where("id = ?", pid).Take(&mediasosial).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this post?
	if uid != mediasosial.UserID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = mediasosial.DeleteAMediaSosial(server.DB, pid, uid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
