package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/orgharoy/healtech/handler"
	"github.com/orgharoy/healtech/middleware"
)

func Routes(app *fiber.App) {

	apiRouter := app.Group("/api")

	apiRouter.Post("/auth/login", handler.Login) //-> Login

	apiRouter.Post("/user", middleware.AuthMiddleware, handler.CreateUser)                                //-> Create User
	apiRouter.Get("/user", middleware.AuthMiddleware, handler.GetActiveUserList)                          // -> Get All Active User List
	apiRouter.Get("/user/maker", middleware.AuthMiddleware, handler.GetUserForMaker)                      // -> Get User List for Checker
	apiRouter.Put("/user/maker/edit/:id", middleware.AuthMiddleware, handler.EditUserFromMakerEnd)        // -> Update User From Maker End
	apiRouter.Put("/user/maker/delete/:id", middleware.AuthMiddleware, handler.DeleteUser)                // -> Delete User From maker End
	apiRouter.Get("/user/checker", middleware.AuthMiddleware, handler.GetUserForChecker)                  // -> Get User List for Checker
	apiRouter.Put("/user/checker/approve/:id", middleware.AuthMiddleware, handler.ApproveUserRecord)      // -> Approve From Checker End
	apiRouter.Put("/user/checker/reject/:id", middleware.AuthMiddleware, handler.SendUserRecordToMaker)   // -> Reject From Checker End
	apiRouter.Put("/user/active/edit/:id", middleware.AuthMiddleware, handler.EditUserFromActiveUserList) // -> Send User Records to Maker end from Activ User's List

	apiRouter.Post("/report", middleware.AuthMiddleware, handler.CreateReport)                              //-> Create Report
	apiRouter.Get("/report/maker", middleware.AuthMiddleware, handler.GetReportListForMaker)                //-> Get Report List For Maker End
	apiRouter.Put("/report/maker/edit/:id", middleware.AuthMiddleware, handler.UpdateReportFromMakerEnd)    //-> Update Report List From Maker End
	apiRouter.Put("/report/maker/delete/:id", middleware.AuthMiddleware, handler.DeleteReport)              // -> Delete Report From Maker End
	apiRouter.Get("/report/checker", middleware.AuthMiddleware, handler.GetReportListForChecker)            //-> Get Report List For Checker End
	apiRouter.Put("/report/checker/approve/:id", middleware.AuthMiddleware, handler.ApproveReportRecord)    // -> Approve From Checker End
	apiRouter.Put("/report/checker/reject/:id", middleware.AuthMiddleware, handler.RejectReportRecord)      // -> Reject From Checker End
	apiRouter.Get("/report", middleware.AuthMiddleware, handler.GetActiveReportList)                        // -> Reject From Checker End
	apiRouter.Put("/report/active/edit/:id", middleware.AuthMiddleware, handler.SendReportToEditFromActive) // -> Send Report Records to Maker end from Activ Report's List

	apiRouter.Get("/report-group", middleware.AuthMiddleware, handler.GetActiveReportGroupList)

	apiRouter.Post("/new-patient-entry", middleware.AuthMiddleware, handler.NewPatientRecord) // -> Handler to Get New Patient Record
	apiRouter.Post("/patient/fetch", middleware.AuthMiddleware, handler.FetchPatientDetails)  // -> Fetch Patient Details

	apiRouter.Get("/bill/:billId", middleware.AuthMiddleware, handler.FetchBillDetails) // -> Fetch bill details with billID
	apiRouter.Put("/bill", middleware.AuthMiddleware, handler.UpdateReportBill) // -> Fetch bill details with billID

	apiRouter.Get("/bill-pdf/download", middleware.AuthMiddleware, handler.GetPDFReport) // -> Fetch bill details with billID

}
