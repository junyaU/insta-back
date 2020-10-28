package controllers

// var globalSessions *session.Manager

// func init() {
// 	globalSessions, _ = session.NewManager("memory", `{"cookieName":"gosessionid", "enableSetCookie,omitempty": true, "gclifetime":3600, "maxLifetime": 3600, "secure": false, "sessionIDHashFunc": "sha1", "sessionIDHashKey": "", "cookieLifeTime": 3600, "providerConfig": ""}`)
// 	go globalSessions.GC()
// }

// func login(w http.ResponceWriter, r *http.Request) {
// 	sess, _ = globalSessions.SessionStart(w, r)
// 	defer sess.SessionRelease(w)
// 	username := sess.Get("username")

// 	if r.Method == "GET" {
// 		t, _ = template.ParseFiles("login.gtpl")
// 	}
// }
