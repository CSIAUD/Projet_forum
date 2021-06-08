package bdd

import (
	"database/sql"
	"fmt"

	"strconv"

	_ "github.com/mattn/go-sqlite3"

	structs "Forum/static/go/structs"
	"time"

	"errors"

	"golang.org/x/crypto/bcrypt"
)

type MyDB struct {
	DB *sql.DB
}

const userRole = 3
const hashCost = 1

//==================================================================================================
func (m MyDB) CreateCommentLike(userid int, commentid int, vote int) bool {
	stmt, err := m.DB.Prepare("INSERT INTO commentLike(user_id, comment_id, vote) values(?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(userid, commentid, vote)
	checkErr(err)

	return true
}
func (m MyDB) UpdateCommentLike(id int, vote int) bool {

	stmt, err := m.DB.Prepare("update commentLike set vote=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(vote, id)
	checkErr(err)

	return true
}
func (m MyDB) DeleteCommentLike(id int) bool {
	stmt, err := m.DB.Prepare("delete from commenLike where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) GetCommentLike(id int) *[]structs.CommentLike {
	rows, err := m.DB.Query("SELECT user_id, comment_id, vote FROM commentLike where id=?", id)
	checkErr(err)
	commentLike := structs.CommentLike{}
	tab := []structs.CommentLike{}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(commentLike.UserId, &commentLike.CommentaireId, &commentLike.Vote)
		checkErr(err)
		tab = append(tab, commentLike)
	}

	return &tab
}

func (m MyDB) CreatePostLike(user_id int, post_id int, vote int) bool {
	stmt, err := m.DB.Prepare("INSERT INTO postLike(user_id, post_id, vote) values(?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(user_id, post_id, vote)
	checkErr(err)

	return true
}
func (m MyDB) UpdatePostLike(id int, vote int) bool {

	stmt, err := m.DB.Prepare("update postLike set vote=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(vote, id)
	checkErr(err)

	return true
}
func (m MyDB) DeletePostLike(id int) bool {
	stmt, err := m.DB.Prepare("delete from posLike where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) GetPostLike(id int) *[]structs.PostLike {
	rows, err := m.DB.Query("SELECT user_id, comment_id, sum(vote) FROM postLike group by post_id where id=?", id)
	checkErr(err)
	postLike := structs.PostLike{}
	tab := []structs.PostLike{}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(postLike.UserId, &postLike.PostId, &postLike.Vote)
		checkErr(err)
		tab = append(tab, postLike)
	}

	return &tab

}

//==================================================================================================
func (m MyDB) AddBadgeUser(user_id int, badge_id int) bool {
	stmt, err := m.DB.Prepare("INSERT INTO badgeUser(user_id, badge_id) values(?,?)")
	checkErr(err)

	_, err = stmt.Exec(user_id, badge_id)
	checkErr(err)

	return true
}
func (m MyDB) DeleteBadgeUser(id int) bool {
	stmt, err := m.DB.Prepare("delete from badgUser where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) GetBadgeUser(user structs.User) *structs.BadgeUser {
	rows, err := m.DB.Query("SELECT badge_id FROM badgeUser where user_id=?", user.Id)
	checkErr(err)
	badgeUser := structs.BadgeUser{}
	badgeUser.User = user
	var badge int

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&badge)
		badgeUser.Badges = append(badgeUser.Badges, (*m.GetBadge(badge)))
		checkErr(err)
	}

	return &badgeUser
}

func (m MyDB) GetBadge(id int) *structs.Badge {
	rows, err := m.DB.Query("SELECT name, image FROM badges where id=?", id)
	checkErr(err)
	temp := structs.Badge{}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&temp.Name, &temp.Image)
		checkErr(err)
	}

	return &temp
}

func (m MyDB) GetAuth(id int) string {
	rows, err := m.DB.Query("SELECT name ROM autorisations where id=?", id)
	checkErr(err)
	var name string

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&name)
		checkErr(err)
	}

	return name
}

func (m MyDB) GetRole(id int) string {
	rows, err := m.DB.Query("SELECT name ROM roles where id=?", id)
	checkErr(err)
	var name string

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&name)
		checkErr(err)
	}

	return name
}

//==================================================================================================
func (m MyDB) Ban(endDate int, raison string, user_id int, bannedBy int) bool {
	stmt, err := m.DB.Prepare("INSERT INTO banList(endDate, raison, bannedBy, userid) values(?,?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(endDate, raison, user_id)
	checkErr(err)

	return true
}
func (m MyDB) BanDef(raison string, user_id int, bannedBy int) bool {
	stmt, err := m.DB.Prepare("INSERT INTO banList(raison, banDef, banneBy, user_id) values(?,?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(raison, 1, bannedBy, user_id)
	checkErr(err)

	return true
}

// func (m MyDB) UnBan(user_id int) bool {
// 	rows, err := m.DB.Query("SELECT id, startDate, raison, banDef, user_id, bannedBy FROM banList where user_id=? ORDER BY startDate desc LIMIT 1", user_id)
// 	checkErr(err)
// 	ban := structsBanList{}
// 	var userid int
// 	var bannedBy int

// 	if rows.Next() {
// 		err = rows.Scan(&ban.Id, &ban.StartDae, &ban.Raison, &ban.BanDef, userid, bannedBy)
// 		ban.BannedBy = *(m.GetUser(banndBy))
// 		ban.User = *(.GetUser(userid))
// 		heckErr(err)
// 	}
// 	rows.Close()

// 	stmt, err := m.DB.Prepare("delete from banList where id=?")
// 	checkErr(err)

// 	_, err = stmtExec((&ban).Id)
//	checkErr(err)

// 	stmt, err = mDB.Prepare("INSERT INTO banList(startDate, raison, user_id, bannedBy) values(?,?,?,?)")
// 	checkErr(err)

// 	_, err = stmtExec((&ban).StartDate, (&ban).Raison, (&ban).UserId, (&ban).BannedBy)
// 	checkErr(err)

// 	eturn true
// }
func (m MyDB) GetBannedUser(user_id int) *[]structs.BanList {
	rows, err := m.DB.Query("SELECT id,startDate,endDate,raison,anDef,bannedBy,user_id FROM banList where user_id=$1 ORDER BY startDate desc", user_id)
	checkErr(err)

	ban := structs.BanList{}
	banList := []structs.BanList{}
	var userid int

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(ban.Id, &ban.StartDate, &ban.EndDate, &ban.Raison, &ban.BanDef, &ban.BannedBy, &userid)
		checkErr(err)
		banList = append(banList, ban)
	}

	return &banList
}

func (m MyDB) GetAllBans() *[]structs.BanList {
	tik := structs.BanList{}
	tab := []structs.BanList{}
	rows, err := m.DB.Query("SELECT t.startdate, t.enddate, t.raison, t.bandef, t.bannedby, user_id FROM banlist t LEFT JOIN users u ON t.user_id=u.id ORDER BY t.StartDate ASC")
	checkErr(err)
	var bannedby int
	var user int

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&tik.StartDate, &tik.EndDate, &tik.Raison, &tik.BanDef, &bannedby, &user)
		checkErr(err)
		tik.BannedBy = (*(m.GetUser(bannedby)))
		tik.User = (*(m.GetUser(user)))
		tab = append(tab, tik)
	}

	return &tab
}

//==================================================================================================
func (m MyDB) CreateUser(username string, mail string, mdp string, avatar string) error {
	rows, err := m.DB.Query("SELECT id FROM users where username like ?", username)
	checkErr(err)

	if rows.Next() {
		return errors.New("error")
	}

	rows, err = m.DB.Query("SELECT id FROM users where mail like ?", mail)
	checkErr(err)

	if rows.Next() {
		return errors.New("error")
	}

	mdp, err = hashMdp(mdp)
	checkErr(err)

	stmt, err := m.DB.Prepare("INSERT INTO users(username, mail, mdp, avatar) values(?,?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(username, mail, mdp, avatar)
	checkErr(err)

	return nil
}
func (m MyDB) UpdasteUser(username string, mail string, avatar string, id int) bool {
	stmt, err := m.DB.Prepare("update users set uername=?, mail=?, avatar=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(username, mail, avatar, id)
	checkErr(err)

	return true
}
func (m MyDB) SetSession(session string, id int) bool {
	stmt, err := m.DB.Prepare("update users set sessionToken=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(session, id)
	checkErr(err)

	return true
}

func (m MyDB) DeleteUser(id int) bool {
	stmt, err := m.DB.Prepare("delete from users where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) GetUser(id int) *structs.User {
	rows, err := m.DB.Query("SELECT id,username,mail,avatar, verified FROM users where id=?", id)
	checkErr(err)
	user := structs.User{}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&user.Id, &user.Username, &user.Mail, &user.Avatar, &user.Verif)
		checkErr(err)
	}
	return &user
}
func (m MyDB) GetUserBySession(token string) *structs.User {
	rows, err := m.DB.Query("SELECT id,username,mail,avatar,role_id,verified FROM users where sessionToken=?", token)
	checkErr(err)
	user := structs.User{}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&user.Id, &user.Username, &user.Mail, &user.Avatar, &user.Role, &user.Verif)
		checkErr(err)
	}
	return &user
}
func (m MyDB) UserExist(mail string) bool {
	fmt.Println(m)
	rows, err := m.DB.Query("SELECT * FROM users where mail=?", mail)

	defer rows.Close()
	fmt.Println("mal2 : ", mail)
	if err != nil {
		return false
	}
	if !rows.Next() {
		return false
	}
	return true
}
func (m MyDB) UserVerified(mail string) bool {
	rows, err := m.DB.Query("SELECT verified FROM users where mail=?", mail)
	checkErr(err)
	var verif int

	checkErr(err)
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&verif)
		checkErr(err)
	}
	if verif == 1 {
		return true
	}
	return false
}

//=================================================================================================================
func (m MyDB) CreatePost(content string, userID int, categorieID int) bool {

	stmt, err := m.DB.Prepare("INSERT INTO posts(content, user_id, categorie_id) values(?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(content, userID, categorieID)
	checkErr(err)

	return true
}
func (m MyDB) UpdatePost(id int, content string, categorieID int, hidden int) bool {
	stmt, err := m.DB.Prepare("update posts set content=?, hidden=?, categorie_id=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(content, hidden, categorieID, id)
	checkErr(err)

	return true
}
func (m MyDB) DeletePost(id int) bool {
	stmt, err := m.DB.Prepare("update posts set hidden=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(1, id)
	checkErr(err)

	// stmt, err := m.DB.Prepare("delete from posts where id=?")
	// checkErr(err)

	// _, err = stmtExec(id)
	// checkErr(err)

	return true
}
func (m MyDB) GetPost(uid int) *structs.Post {
	post := structs.Post{}

	rows, err := m.DB.Query("SELECT p.id, p.content, p.date, p.categorie_id, p.hidden, p.user_id, u.username, u.avatar FROM posts p LEFT JOIN users u ON u.id = p.user_id WHERE p.id=?", uid)
	checkErr(err)
	var date int
	var cat int

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&post.Id, &post.Content, &date, &cat, &post.Hidden, &post.User.Id, &post.User.Username, &post.User.Avatar)
		post.Date = m.DateConversion(date)
		post.Categorie = m.GetCategory(cat)
		checkErr(err)
	}
	return &post
}
func (m MyDB) GetNbPost(limit int, offset int) *[]structs.Post {
	offset = offset * limit
	rows, err := m.DB.Query("SELECT p.id, p.content, p.date, p.categorie_id, p.hidden, p.user_id, u.username, u.avatar FROM posts p LEFT JOIN users u ON u.id = p.user_id WHERE hidden!=1 ORDER BY date asc LIMIT ? OFFSET ?", limit, offset)
	checkErr(err)

	post := structs.Post{}
	tab := []structs.Post{}
	var cat int
	var date int

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&post.Id, &post.Content, &date, &cat, &post.Hidden, &post.User.Id, &post.User.Username, &post.User.Avatar)
		checkErr(err)
		post.Categorie = m.GetCategory(cat)
		post.Date = m.DateConversion(date)
		tab = append(tab, post)
	}

	return &tab
}

func (m MyDB) CreateComment(content string, userId int, postId int, commentId int) bool {
	to := ""
	if commentId == 0 {
		to = "post_id"
	} else if postId == 0 {
		to = "commentaire_id"
	}
	stmt, err := m.DB.Prepare("INSERT INTO commentaires(content, user_id, " + to + ") values(?,?,?)")
	checkErr(err)

	if commentId == 0 {
		_, err = stmt.Exec(content, userId, postId)
	} else if postId == 0 {
		_, err = stmt.Exec(content, userId, commentId)
	}
	checkErr(err)

	return true
}
func (m MyDB) UpdateComment(id int, content string) bool {

	stmt, err := m.DB.Prepare("update commentaires set content=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(content, id)
	checkErr(err)

	return true
}
func (m MyDB) DeleteComment(id int) bool {
	stmt, err := m.DB.Prepare("delete from comentaires where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) GetComment(uid int) *[]structs.Commentaire {
	comment := structs.Commentaire{}
	rows, err := m.DB.Query("SELECT c.id, c.content, c.date, c.user_id, u.username, u.avatar FROM commentaires c LEFT JOIN users u ON u.id = c.user_id WHERE c.post_id=?", uid)
	tab := []structs.Commentaire{}
	var date int

	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&comment.Id, &comment.Content, &date, &comment.User.Id, &comment.User.Username, &comment.User.Avatar)
		comment.Date = m.DateConversion((date))
		checkErr(err)
		tab = append(tab, comment)
	}

	return &tab
}

//=============================================================================================
func (m MyDB) CreateCategory(name string) bool {
	stmt, err := m.DB.Prepare("INSERT INTO categories(name) values(?)")
	checkErr(err)

	_, err = stmt.Exec(name)
	checkErr(err)

	return true
}
func (m MyDB) UpdateCategory(id int, name string) bool {

	stmt, err := m.DB.Prepare("update categories set name=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(name, id)
	checkErr(err)

	return true
}
func (m MyDB) DeleteCategory(id int) bool {
	stmt, err := m.DB.Prepare("delete from catgories where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) GetCategory(id int) string {
	rows, err := m.DB.Query("SELECT name FROM categories where id=?", id)
	checkErr(err)
	var name string

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&name)
		checkErr(err)
	}

	return name
}

//===========================================================================

func (m MyDB) CreateTicket(id int, content string, categorieId int) bool {
	stmt, err := m.DB.Prepare("INSERT INTO tickets(content, user_id, categorie_id) values(?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(id, content, categorieId)
	checkErr(err)

	return true
}
func (m MyDB) OpenTicket(id int) bool {

	stmt, err := m.DB.Prepare("update categories set etat=1 where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) CloseTicket(id int) bool {

	stmt, err := m.DB.Prepare("update categories set etat=2 where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) GetAllTickt() *[]structs.Ticket {
	tik := structs.Ticket{}
	tab := []structs.Ticket{}
	rows, err := m.DB.Query("SELECT t.id, t.content, t.date, t.etat, t.categorie_id, t.openBy, u.id user FROM tickets t LEFT JOIN users u ON t.user_id=u.id ORDER BY date ASC")
	checkErr(err)
	var openBy int
	var user int

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&tik.Id, &tik.Content, &tik.Date, &tik.Etat, &tik.Categorie, &openBy, &user)
		checkErr(err)
		tik.User = *(m.GetUser(user))
		tik.OpenBy = *(m.GetUser(openBy))
		tab = append(tab, tik)
	}

	return &tab
}
func (m MyDB) GetTicket(id int) *structs.Ticket {
	tik := structs.Ticket{}
	rows, err := m.DB.Query("SELECT t.id, t.content, t.date, t.etat, t.categorie_id, t.openBy, u.id user FROM tickets t LEFT JOIN users u ON t.user_id=u.id WHERE t.id=? ORDER BY date ASC", id)
	checkErr(err)
	var openBy int
	var user int

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&tik.Id, &tik.Content, &tik.Date, &tik.Etat, &tik.Categorie, &openBy, &user)
		checkErr(err)
		tik.User = (*m.GetUser(user))
		tik.OpenBy = (*m.GetUser(openBy))
	}

	return &tik
}

//=============================================================================
func hashMdp(mdp string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(mdp), hashCost)
	return string(bytes), err
}
func (m MyDB) CompareMdp(password string, mail string) (int, error) {
	id := 0
	rows, err := m.DB.Query("SELECT id, mdp FROM users where mail=?", mail)
	defer rows.Close()
	if err != nil {
		return 0, err
	}
	var mdp string

	if !rows.Next() {
		return id, errors.New("error")
	}
	err = rows.Scan(&id, &mdp)
	checkErr(err)

	err = bcrypt.CompareHashAndPassword([]byte(mdp), []byte(password))
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (m MyDB) updateMdp(old string, mdp string, mail string) bool {
	_, err := m.CompareMdp(old, mail)
	if err != nil {
		return false
	}
	stmt, err := m.DB.Prepare("update users set mdp=? where mail=?")
	checkErr(err)

	_, err = stmt.Exec(mdp, mail)
	checkErr(err)
	return true
}
func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func (m MyDB) DateConversion(date int) string {

	now := int(time.Now().Unix()) + 2*3600
	diff := (now - date) // Calcul du nombre de secondes entre la date de création et maintenant
	temp := ""

	if diff < 60 {
		temp += "moins d'une minute"
	} else {
		diff /= 60 // n passe en Minutes
		if diff < 60 {
			temp += strconv.Itoa(diff) + " minutes"
		} else {
			diff /= 60 // n passe en heures
			if diff == 1 {
				temp += "1 heure"
			} else if diff < 24 {
				temp += strconv.Itoa(diff) + " heures"
			} else {
				diff /= 24 // n passe en jours
				if diff == 1 {
					temp += "1 jour"
				} else if diff < 30 {
					temp += strconv.Itoa(diff) + " jours"
				} else {
					diff /= 30 // n passe en mois
					if diff < 12 {
						temp += strconv.Itoa(diff) + " mois"
					} else {
						diff /= 12 // n passe en années
						if diff == 1 {
							temp += "1 an"
						} else {
							temp += strconv.Itoa(diff) + " ans"
						}
					}
				}

			}
		}
	}

	// if diff < 60 {
	// 	fmt.Println(diff)
	// 	temp +=strconv.Itoa(diff) + " Min.s"
	// } else {
	// 	diff /= 60
	// 	if diff < 60 {
	// 		fmt.Println(diff)
	// 		temp +=strconv.Itoa(diff) + " H"
	// 	} else {
	// 		diff /= 60
	// 		if diff < 24 {
	// 			fmt.Println(diff)
	// 			temp +=strconv.Itoa(diff) + " Jour.s"
	// 		} else {
	// 			fmt.Printl(diff)
	// 			diff /= 24
	// 			if diff < 30 {
	// 				fmt.Println(diff)
	// 				temp +=strconv.Itoa(diff) + " Mois"
	// 			} else {
	// 				diff /= 30
	// 				diff /= 12
	// 				fmt.Println(diff)
	// 				emp += strconv.Itoa(diff) + " An.s"
	//
	//
	//
	// }

	return temp
}
