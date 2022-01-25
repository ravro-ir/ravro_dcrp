package main

import (
	"fmt"
	ptime "github.com/yaa110/go-persian-calendar"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Pdf struct {
	report    Report
	judge     Judgment
	amendment Amendment
}

func main() {
	var (
		templatePath string
		outputPath   string
	)
	publicMessage := "شرح داده نشد است"
	defer func() {
		if err := recover(); err != nil {
			log.Println("Error, please check your input encrypt file Or report the issue in the github.")
		}
	}()
	r := NewRequestPdf("")
	pt := ptime.Now()
	AddDir("decrypt")
	if runtime.GOOS == "windows" {
		templatePath = "template\\sample.html"
		outputPath = "decrypt\\reports.pdf"
	} else {
		templatePath = "template/sample.html"
		outputPath = "decrypt/reports.pdf"
	}
	fmt.Println("[++++] Starting for decrypting . . . ")
	judge, err := DcrptJudgment()
	if err != nil {
		log.Fatalln(err)
	}
	report, err := DcrptReport()
	if err != nil {
		log.Fatalln(err)
	}
	amendment, err := DcrptAmendment()
	if err != nil {
		log.Fatalln(err)
	}
	moreInfo := strings.Join(amendment[:], ",")
	if moreInfo == "" {
		moreInfo = publicMessage
	}
	dateTo := strconv.Itoa(pt.Year()) + "/" + strconv.Itoa(int(pt.Month())) + "/" + strconv.Itoa(pt.Day())
	pdf := Pdf{judge: judge, report: report}
	outputPath = strings.Replace(outputPath, "reports", report.CompanyUsername+"__"+report.Slug+"__"+report.HunterUsername, 1)
	if pdf.report.Reproduce == "" {
		pdf.report.Reproduce = publicMessage
	}
	if pdf.judge.Description == "" {
		pdf.judge.Description = publicMessage
	}
	fmt.Println("[++++] decrypted successfully ")
	dateSubmit := pdf.report.SubmissionDate

	dateSubmited := strings.Split(dateSubmit, " ")
	dateReport := strings.Split(string(dateSubmited[0]), "-")
	year, _ := strconv.Atoi(dateReport[0])
	month, _ := strconv.Atoi(dateReport[1])
	day, _ := strconv.Atoi(dateReport[2])
	var t time.Time = time.Date(year, time.Month(month), day, 12, 1, 1, 0, ptime.Iran())
	pt = ptime.New(t)
	dataFrom := pt.Format("yyyy/MM/dd")
	fmt.Println("[++++] Starting report to pdf . . . ")
	templateData := struct {
		Title           string
		Description     string
		PoC             string
		DateFrom        string
		CVSS            string
		Reproduce       string
		Hunter          string
		ReportID        string
		Amount          int
		Score           int
		JudgeInfo       string
		DateTo          string
		MoreInfo        string
		CompanyUserName string
	}{
		Title:           pdf.report.Title,
		PoC:             pdf.report.Description,
		CVSS:            pdf.judge.Cvss.Value,
		Reproduce:       pdf.report.Reproduce,
		DateFrom:        dataFrom,
		Hunter:          pdf.report.HunterUsername,
		ReportID:        pdf.report.Slug,
		Amount:          pdf.judge.Reward,
		JudgeInfo:       pdf.judge.Description,
		DateTo:          dateTo,
		MoreInfo:        moreInfo,
		CompanyUserName: pdf.report.CompanyUsername,
	}
	if err := r.ParseTemplate(templatePath, templateData); err == nil {
		_, _ = r.GeneratePDF(outputPath)
		fmt.Println("[++++] pdf generated successfully")
	} else {
		fmt.Println(err)
	}
}
