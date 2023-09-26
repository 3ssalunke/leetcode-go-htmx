package main

// func (server *Server) TestAndVerfiyUserCode(w http.ResponseWriter, r *http.Request) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(300)*time.Second)
// 	defer cancel()

// 	var requestData CodeExecuteRequest
// 	err := json.NewDecoder(r.Body).Decode(&requestData)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	problems, err := controllers.GetProblemDetailsByProblemID(ctx, server.Db, requestData.ProblemID)
// 	if err != nil {
// 		log.Printf("failed to fetch details for problem with id %s - %v", requestData.ProblemID, err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	for i, value := range problems[0].TestCaseList {
// 		err = util.WriteCodeInExecutionFile(requestData.Lang, requestData.TypedCode, problems[0].SolutionName, value)

// 		if err != nil {
// 			log.Printf("failed to write to the file - %v", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		output, err := util.ExecuteCode(ctx, requestData.Lang)
// 		if err != nil {
// 			log.Printf("failed to execute the code - %v", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		output = strings.ReplaceAll(output, " ", "")

// 		if !reflect.DeepEqual([]byte(output), []byte(problems[0].TestCaseAnswers[i])) {
// 			log.Printf("test case %d failed", i)
// 		} else {
// 			log.Printf("test case %d passed", i)
// 		}
// 	}

// 	w.WriteHeader(http.StatusOK)
// }
