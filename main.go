package main

import (
	"email_action/jobs"
	"email_action/logging"
	"email_action/mail"
	"email_action/routers"
	_ "email_action/routers"
	"email_action/store"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	cronv3 "github.com/robfig/cron/v3"
)

var (
	log = logging.NewZapLogger()
)

func main() {
	err := beego.LoadAppConfig("ini", "conf/app.conf")
	if err != nil {
		log.Errorf("failed to read config. ABORT!!")
		return
	}

	env := beego.AppConfig.DefaultString("env", "dev")
	port := beego.AppConfig.DefaultInt("httpport", 80)
	jwtSecret := beego.AppConfig.String("jwtecret")

	store := store.NewStore(env, jwtSecret)
	mailer := mail.NewSesAdapter()
	routers.InitRouters(store, mailer)

	cron := cronv3.New()
	sendSubscriptionEmailsJob := jobs.NewSendSubscriptionEmailsJob(store, mailer)
	clearFreePlanEmailQuotaJob := jobs.NewClearFreePlanEmailQuotaJob(store, mailer)

	cronJobs := []jobs.JobInterface{sendSubscriptionEmailsJob, clearFreePlanEmailQuotaJob}
	cronInterval1 := beego.AppConfig.DefaultString("cron_interval_send_email", "0 10 * * *") // //https://crontab.guru/
	cronInterval2 := beego.AppConfig.DefaultString("cron_interval_clear_quota", "0 6 * * *")
	cronIntervals := []string{cronInterval1, cronInterval2}
	for i, job := range cronJobs {
		job := job
		cron.AddFunc(cronIntervals[i], func() { job.Run() })
	}
	cron.Start()
	log.Infof("server is running in %s at port %d", env, port)

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin", "content-type", "Access-Control-Allow-Origin", "authorization"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
	}))

	beego.Run()
}
