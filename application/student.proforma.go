package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spo-iitk/ras-backend/rc"
	"github.com/spo-iitk/ras-backend/util"
)

func getProformasForStudentHandler(ctx *gin.Context) {
	rid, err := util.ParseUint(ctx.Param("rid"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var jps []Proforma

	err = fetchProformasForStudent(ctx, rid, &jps)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, jps)
}

func getProformasForEligibleStudentHandler(ctx *gin.Context) {
	rid, err := util.ParseUint(ctx.Param("rid"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sid := getStudentRCID(ctx)
	var student rc.StudentRecruitmentCycle

	err = rc.FetchStudent(ctx, sid, &student)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var jps []Proforma

	err = fetchProformaForEligibleStudent(ctx, rid, &student, &jps)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, jps)
}

func getProformaForStudentHandler(ctx *gin.Context) {
	pid, err := util.ParseUint(ctx.Param("pid"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch the student's RC record for eligibility check
	sid := getStudentRCID(ctx)
	if sid == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Student RC ID not found"})
		return
	}

	var student rc.StudentRecruitmentCycle
	err = rc.FetchStudent(ctx, sid, &student)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var jp Proforma

	err = fetchProformaForStudent(ctx, pid, &jp)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check department-wise eligibility before returning proforma details
	primaryID := int(student.ProgramDepartmentID)
	secondaryID := int(student.SecondaryProgramDepartmentID)
	eligibility := jp.Eligibility

	primaryEligible := primaryID > 0 && primaryID < len(eligibility) && eligibility[primaryID] == '1'
	secondaryEligible := secondaryID > 0 && secondaryID < len(eligibility) && eligibility[secondaryID] == '1'

	if !primaryEligible && !secondaryEligible {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Your department is not eligible for this opening"})
		return
	}

	ctx.JSON(http.StatusOK, jp)
}
