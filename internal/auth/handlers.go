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
