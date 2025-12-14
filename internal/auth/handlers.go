package auth

/*

// @Summary	Download images
// @Tags		Task
// @Produce	json
// @Param		request	body		model.Urls	true	"Urls"
// @Success	200		{object}	model.DownloadImages
// @Failure	400		{object}	model.APIError
// @Failure	408     {object}	model.APIError
// @Failure	500		{object}	model.APIError
// @Router		/download/images [post]
func DownloadUrls(ctx *gin.Context) {
	urls := models.Urls{}
	err := ctx.ShouldBindJSON(&urls)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "INVALID_REQUEST", err))
		return
	}

	resultCh := make(chan models.DownloadImages, len(urls.Urls))
	for _, url := range urls.Urls {
		go downLoadImage(url, resultCh)
	}

	timeout := time.After(10 * time.Second)

	results := make([]models.DownloadImages, 0, len(urls.Urls))
	for i := 0; i < len(urls.Urls); i++ {
		select {
		case result := <-resultCh:
			results = append(results, result)
		case <-timeout:
			ctx.JSON(http.StatusRequestTimeout, gin.H{
				"error":   "Timeout waiting for downloads",
				"results": results,
			})
			return
		}

	}

	apperrors.GetAPIError(ctx, gin.H{"results": results}, http.StatusOK, nil)
}

// @Summary	Query Task (server side filtering & paging)
// @Tags		Task
// @Produce	json
// @Param		request	body		model.QueryTasksRequest	true	"QueryTasksRequest"
// @Success	200		{object}	model.QueryTasksResponse
// @Failure	400		{object}	model.APIError
// @Failure	500		{object}	model.APIError
// @Router		/tasks/query [post]
func QueryTasksV2(ctx *gin.Context) {
	queryRequest := models.QueryTasksRequest{}
	err := ctx.ShouldBindJSON(&queryRequest)
	if err != nil {
		apperrors.GetAPIError(ctx, nil, 0, apperrors.APIError{}.New(http.StatusBadRequest, "INVALID_REQUEST", err))
		return
	}

	var filtered []models.Task
	for _, task := range models.LibTasks {
		if strings.Contains(strings.ToLower(task.Title), strings.ToLower(queryRequest.Search)) {
			filtered = append(filtered, *task)
		}
	}

	totalRows := len(filtered)
	if totalRows == 0 {
		apperrors.GetAPIError(ctx, models.QueryTasksResponse{TotalResults: totalRows, TotalPages: 0, Tasks: []models.Task{}}, http.StatusOK, nil)
		return
	}

	totalPages := int(math.Ceil(float64(totalRows) / float64(queryRequest.Paging.PageSize)))

	startPage := (queryRequest.Paging.Page - 1) * queryRequest.Paging.PageSize
	endPage := startPage + queryRequest.Paging.PageSize

	switch {
	case startPage > totalRows:
		startPage = totalRows
	case endPage > totalRows:
		endPage = totalRows
	}

	pagedTasks := filtered[startPage:endPage]

	response := models.QueryTasksResponse{
		TotalResults: totalRows,
		TotalPages:   totalPages,
		Tasks:        pagedTasks,
	}
	apperrors.GetAPIError(ctx, response, http.StatusOK, nil)
}
*/

/*func downLoadImage(url string, ch chan<- models.DownloadImages) {
	client := http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		ch <- models.DownloadImages{
			URL:     url,
			Success: false,
			Error:   fmt.Errorf("HTTP error: %v", err),
		}
		return
	}
	defer resp.Body.Close()

	ch <- models.DownloadImages{
		URL:     url,
		Success: true,
		Error:   nil,
	}
}*/
