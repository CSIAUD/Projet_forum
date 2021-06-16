package bdd

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	structs "Forum/static/go/structs"
	"time"

	"errors"

	"sync"

	"golang.org/x/crypto/bcrypt"
)

type MyDB struct {
	DB *sql.DB
}

const userRole = 3
const hashCost = 1

// LIKES ==================================================================================================
func (m MyDB) CreateCommentLike(userid int, commentid int, vote int) bool { //Pouvoir liker un commentaire
	stmt, err := m.DB.Prepare("INSERT INTO commentLike(user_id, comment_id, vote) values(?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(userid, commentid, vote)
	checkErr(err)

	return true
}
func (m MyDB) UpdateCommentLike(id int, vote int) bool { // Pouvoir changer son like en dislike et inversement

	stmt, err := m.DB.Prepare("update commentLike set vote=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(vote, id)
	checkErr(err)

	return true
}
func (m MyDB) DeleteCommentLike(id int) bool { // Pouvoir retirer complètement son vote
	stmt, err := m.DB.Prepare("delete from commenLike where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) GetCommentLike(id int) *structs.Like { // Récupèrer les données des votes de commentaires grâce aux structs
	rows, err := m.DB.Query("SELECT sum(vote), count(vote) FROM commentLike where comment_id=? group by comment_id", id)
	checkErr(err)
	note := structs.Like{}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&note.Note, &note.Total)
		checkErr(err)
	}

	return &note
}

func (m MyDB) CreatePostLike(user_id int, post_id int, vote int) bool { // Pouvoir liker un post
	stmt, err := m.DB.Prepare("INSERT INTO postLike(user_id, post_id, vote) values(?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(user_id, post_id, vote)
	if err != nil {
		if m.GetPostLikeUser(user_id, post_id) != vote {
			m.UpdatePostLike(user_id, post_id, vote)
		} else {
			m.DeletePostLike(user_id, post_id)
		}
	}

	return true
}
func (m MyDB) UpdatePostLike(user int, post int, vote int) bool { // Pouvoir modifier la valeur du vote du post

	stmt, err := m.DB.Prepare("update postLike set vote=? where user_id=? and post_id=?")
	checkErr(err)

	_, err = stmt.Exec(vote, user, post)
	checkErr(err)

	return true
}
func (m MyDB) DeletePostLike(user int, post int) bool { // Pouvoir supprimer le vote du post
	stmt, err := m.DB.Prepare("delete from postLike where user_id=? and post_id=?")
	checkErr(err)

	_, err = stmt.Exec(user, post)
	checkErr(err)

	return true
}
func (m MyDB) GetPostLikeUser(user int, post int) int { // Récupérer les données du vote du post d'un utilisateur en particulier (nous) grâce aux structs
	rows, err := m.DB.Query("SELECT vote FROM postLike where user_id=? and post_id=?", user, post)
	checkErr(err)
	var count int

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&count)
		checkErr(err)
	}

	return count
}
func (m MyDB) GetPostLike(id int) *structs.Like { // Récupérer les données du vote du post grâce aux structs
	rows, err := m.DB.Query("SELECT sum(vote), count(vote) FROM postLike where post_id=? group by post_id", id)
	checkErr(err)
	note := structs.Like{}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&note.Note, &note.Total)
		checkErr(err)
	}

	return &note
}
func (m MyDB) GetLike(id int) (string, string) { // Récupérer les données du vote du post grâce aux structs
	rows, err := m.DB.Query("SELECT sum(vote), count(vote) FROM postLike where post_id=? group by post_id", id)
	checkErr(err)
	var note int
	var total int
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&note, &total)
		checkErr(err)
	}

	return strconv.Itoa(note), strconv.Itoa(total)
}

// BADGES / ROLES ==================================================================================================
func (m MyDB) AddBadgeUser(user_id int, badge_id int) bool { // Pouvoir ajouter un badge à un utilisateur
	stmt, err := m.DB.Prepare("INSERT INTO badgeUser(user_id, badge_id) values(?,?)")
	checkErr(err)

	_, err = stmt.Exec(user_id, badge_id)
	checkErr(err)

	return true
}
func (m MyDB) DeleteBadgeUser(id int) bool { // Pouvoir supprimer un badge à un utilisateur
	stmt, err := m.DB.Prepare("delete from badgUser where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) GetBadgeUser(user structs.User) *structs.BadgeUser { // Récupérer les données des badges d'un utilisateur en particulier (nous) grâce aux structs
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

func (m MyDB) GetBadge(id int) *structs.Badge { // Récupérer les données des badges grâce aux structs
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

func (m MyDB) GetRole(id int) string { // Récupérer les données des rôles grâce aux structs
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

func (m MyDB) UpdateRole(role int, id int) bool { // Pouvoir modifier le rôle d'un utilisateur (user, modo, admin)
	stmt, err := m.DB.Prepare("update users set role_id=? where id=? and role_id <> ?")
	checkErr(err)

	_, err = stmt.Exec(role, id, role)
	checkErr(err)
	return true
}

// BANNISSEMENTS ==================================================================================================
func (m MyDB) Ban(endDate int, raison string, user_id int, bannedBy int) bool { // Pouvoir insérer un ban dans les données d'un utilisateur
	stmt, err := m.DB.Prepare("INSERT INTO banList(endDate, raison, bannedBy, userid) values(?,?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(endDate, raison, user_id)
	checkErr(err)

	return true
}
func (m MyDB) BanDef(raison string, user_id int, bannedBy int) bool { // Pouvoir bannir définitivement un utilisateur
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
//
// 	if rows.Next() {
// 		err = rows.Scan(&ban.Id, &ban.StartDae, &ban.Raison, &ban.BanDef, userid, bannedBy)
// 		ban.BannedBy = *(m.GetUser(banndBy))
// 		ban.User = *(.GetUser(userid))
// 		heckErr(err)
// 	}
// 	rows.Close()
//
// 	stmt, err := m.DB.Prepare("delete from banList where id=?")
// 	checkErr(err)
//
// 	_, err = stmtExec((&ban).Id)
//	checkErr(err)
//
// 	stmt, err = mDB.Prepare("INSERT INTO banList(startDate, raison, user_id, bannedBy) values(?,?,?,?)")
// 	checkErr(err)
//
// 	_, err = stmtExec((&ban).StartDate, (&ban).Raison, (&ban).UserId, (&ban).BannedBy)
// 	checkErr(err)
//
// 	eturn true
// }
func (m MyDB) GetBannedUser(user_id int) *[]structs.BanList { // Récupérer les données du ban d'un utilisateur en particulier grâce aux structs
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
func (m MyDB) GetAllBans() *[]structs.BanList { // Récupérer les données du ban de tous les utilisateurs bannis grâce aux structs
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

// USERS ==================================================================================================
func (m MyDB) CreateUser(username string, mail string, mdp string) error { // Pouvoir s'inscrire au site

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

	stmt, err := m.DB.Prepare("INSERT INTO users(username, mail, mdp) values(?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(username, mail, mdp)
	checkErr(err)

	return nil
}
func (m MyDB) UpdateUser(username string, mail string, avatar string, id int) bool { // Pouvoir modifier ses informations personnelles
	stmt, err := m.DB.Prepare("update users set username=?, mail=?, avatar=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(username, mail, avatar, id)
	checkErr(err)

	return true
}
func (m MyDB) SetSession(session string, id int) bool { //mise à jour de la valeur du sessionToken dans la bdd du User
	stmt, err := m.DB.Prepare("update users set sessionToken=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(session, id)
	checkErr(err)

	return true
}

func (m MyDB) DeleteUser(id int) bool { // Pouvoir supprimer son compte / se faire supprimer son compte
	stmt, err := m.DB.Prepare("delete from postLike where post_id in(select id posts from posts where user_id=?)")
	checkErr(err)
	_, err = stmt.Exec(id)
	checkErr(err)

	stmt, err = m.DB.Prepare("delete from commentaires where post_id in(select id posts from posts where user_id=?)")
	checkErr(err)
	_, err = stmt.Exec(id)
	checkErr(err)

	stmt, err = m.DB.Prepare("delete from posts where user_id=?")
	checkErr(err)
	_, err = stmt.Exec(id)
	checkErr(err)

	stmt, err = m.DB.Prepare("delete from users where id=?")
	checkErr(err)
	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}

func (m MyDB) DeleteMyUser(id int) bool { // Pouvoir supprimer son compte / se faire supprimer son compte
	stmt, err := m.DB.Prepare("update users set username=deleted, mail=deleted, mdp=deleted where id=?")
	checkErr(err)
	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}

func (m MyDB) GetUser(id int) *structs.User { // Récupérer les données de l'utilisateur grâce aux structs
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
func (m MyDB) GetUserByRole(role int) *[]structs.User { // Récupérer les données de l'utilisateur grâce à son rôle
	tik := structs.User{}
	tab := []structs.User{}
	rows, err := m.DB.Query("SELECT id, username FROM users WHERE role_id=?", role)
	checkErr(err)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&tik.Id, &tik.Username)
		checkErr(err)
		tab = append(tab, tik)
	}

	return &tab
}
func (m MyDB) GetUserByName(username string) *structs.User { // Récupérer les données d'un utilisateur par son nom
	rows, err := m.DB.Query("SELECT id,username,mail, verified FROM users where username=?", username)
	checkErr(err)
	user := structs.User{}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&user.Id, &user.Username, &user.Mail, &user.Verif)
		checkErr(err)
	}
	return &user
}
func (m MyDB) GetUserByMail(mail string) *structs.User { // Récupérer les données d'un utilisateur grâce à son mail
	rows, err := m.DB.Query("SELECT id, mail FROM users where mail=?", mail)
	checkErr(err)
	user := structs.User{}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&user.Id, &user.Mail)
		checkErr(err)
	}
	fmt.Println("user in db :", user)
	return &user
}
func (m MyDB) GetUserModo() *[]structs.UserCat { // Récupérer les données de l'utilisateur grâce à son rôle
	tab := []structs.UserCat{}
	user := structs.UserCat{}
	rows, err := m.DB.Query("SELECT id, username FROM users WHERE role_id=2")
	checkErr(err)

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Username)
		checkErr(err)
		tab = append(tab, user)
	}

	wait := sync.WaitGroup{}
	lock := sync.Mutex{}
	for i := 0; i < len(tab); i++ {
		wait.Add(1)
		go func(index int) {
			id := tab[index].Id
			rows, err := m.DB.Query("select categorie_id from userCat where user_id=?", id)
			checkErr(err)
			var catTab []int
			var cat int

			defer rows.Close()
			for rows.Next() {
				err = rows.Scan(&cat)
				checkErr(err)
				catTab = append(catTab, cat)
			}
			lock.Lock()
			tab[index].Categories = catTab
			lock.Unlock()
			wait.Done()
		}(i)
	}
	wait.Wait()
	fmt.Println(tab)
	return &tab
}
func (m MyDB) GetUserPost(uid int) *[]structs.Post { // Récupérer les posts d'un utilisateur en particulier
	fmt.Println("USER ID FOR POSTS:", uid)
	post := structs.Post{}
	tab := []structs.Post{}

	rows, err := m.DB.Query("SELECT p.id, p.content, p.date, p.categorie_id, p.hidden, p.user_id, u.username, u.avatar FROM posts p LEFT JOIN users u ON u.id = p.user_id WHERE p.user_id=?", uid)
	checkErr(err)
	var date int
	var cat int

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&post.Id, &post.Content, &date, &cat, &post.Hidden, &post.User.Id, &post.User.Username, &post.User.Avatar)
		post.Date = m.DateConversion(date)
		post.Categorie = m.GetCategory(cat)
		tab = append(tab, post)
		checkErr(err)
	}
	return &tab
}
func (m MyDB) GetUserComment(uid int) *[]structs.Commentaire { // Récupérer les commentaires d'un utilisateur en particulier
	comment := structs.Commentaire{}
	rows, err := m.DB.Query("SELECT c.id, c.content, c.date, c.user_id, u.username, u.avatar FROM commentaires c LEFT JOIN users u ON u.id = c.user_id WHERE c.user_id=?", uid)
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
func (m MyDB) GetTheUserPost(user structs.User) *structs.Posts { // Récupérer les posts d'un utilisateur en particulier (nous)
	rows, err := m.DB.Query("SELECT id, content, date, categorie_id, hidden FROM posts where user_id=?", user.Id)
	checkErr(err)
	postUser := structs.Posts{}
	post := structs.Post{}
	postUser.User = user
	var date int
	var cat int

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&post.Content, &date, &cat, &post.Hidden)
		post.Date = m.DateConversion(date)
		post.Categorie = m.GetCategory(cat)
		postUser.Posts = append(postUser.Posts, post)
		checkErr(err)
	}

	return &postUser
}

func (m MyDB) GetUserBySession(token string) *structs.User { // Récupérer les données de l'utilisateur grâce à son cookie Session Token
	user := structs.User{}
	if token == "0" {
		user.Id = 0
		user.Username = ""
		user.Mail = ""
		user.Avatar = ""
		user.SessionToken = ""
		user.Role = 0
		user.Verif = 0
		return &user
	}
	rows, err := m.DB.Query("SELECT id,username,mail,avatar,role_id,verified FROM users where sessionToken=?", token)
	checkErr(err)

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&user.Id, &user.Username, &user.Mail, &user.Avatar, &user.Role, &user.Verif)
		checkErr(err)
	}
	fmt.Printf("Session => %s\nUser : ", token)
	fmt.Println(user)
	return &user
}
func (m MyDB) UserExist(mail string) bool { // Pouvoir vérifier si un utilisateur existe grâce à son mail pendant l'inscription
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
func (m MyDB) UserVerified(mail string) bool { // Pouvoir finaliser l'inscription en validant son mail
	stmt, err := m.DB.Prepare("update users set verified=1 where mail=?")
	checkErr(err)

	_, err = stmt.Exec(mail)
	checkErr(err)

	return true
}

func (m MyDB) CreateUserToken(userId int, token string) bool {
	stmt, err := m.DB.Prepare("INSERT INTO userToken(user_id, token) values(?,?)")
	checkErr(err)

	_, err = stmt.Exec(userId, token)
	checkErr(err)
	if err == nil {
		fmt.Println("UPDATE USERTOKEN")
	}

	return true
}

func (m MyDB) GetUserToken(token string) *structs.UserToken {
	rows, err := m.DB.Query("SELECT user_id, token FROM userToken where token=?", token)
	checkErr(err)
	userToken := structs.UserToken{}

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&userToken.Userid, &userToken.Token)
		checkErr(err)
	}
	return &userToken
}

// POSTS / COMMENTAIRES =================================================================================================================
func (m MyDB) CreatePost(content string, userID int, categorieID int) bool { // Pouvoir créer un post dans une catégorie spécifique

	stmt, err := m.DB.Prepare("INSERT INTO posts(content, user_id, categorie_id) values(?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(content, userID, categorieID)
	checkErr(err)

	return true
}
func (m MyDB) UpdatePost(id int, content string, categorieID int, hidden int) bool { // Pouvoir modifier son post
	stmt, err := m.DB.Prepare("update posts set content=?, hidden=?, categorie_id=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(content, hidden, categorieID, id)
	checkErr(err)

	return true
}
func (m MyDB) DeletePost(id int) bool { // Pouvoir supprimer son post ou se le faire supprimer
	stmt, err := m.DB.Prepare("update posts set hidden=1 where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	// stmt, err := m.DB.Prepare("delete from posts where id=?")
	// checkErr(err)

	// _, err = stmtExec(id)
	// checkErr(err)

	return true
}
func (m MyDB) GetPost(uid int) *structs.Post { // Récupérer les données d'un post grâce aux structs
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
		post.Likes = (*m.GetPostLike(post.Id))
		checkErr(err)
	}
	return &post
}
func (m MyDB) GetNbPost(limit int, offset int) *[]structs.Post { // Pouvoir afficher un certain nombre de posts dans l'index
	offset = offset * limit
	rows, err := m.DB.Query("SELECT p.id, p.content, p.date, p.categorie_id, p.hidden, p.user_id, u.username, u.avatar FROM posts p LEFT JOIN users u ON u.id = p.user_id WHERE hidden!=1 ORDER BY date desc LIMIT ? OFFSET ?", limit, offset)
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
		post.Likes = (*m.GetPostLike(post.Id))
		tab = append(tab, post)
	}

	return &tab
}

func (m MyDB) CreateComment(content string, userId int, postId int) bool { // Pouvoir créer un commentaire

	stmt, err := m.DB.Prepare("INSERT INTO commentaires(content, user_id, post_id) values(?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(content, userId, postId)
	checkErr(err)

	return true
}
func (m MyDB) UpdateComment(id int, content string) bool { // Pouvoir modifier son commentaire

	stmt, err := m.DB.Prepare("update commentaires set content=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(content, id)
	checkErr(err)

	return true
}
func (m MyDB) DeleteComment(id int) bool { // Pouvoir supprimer son commentaire
	stmt, err := m.DB.Prepare("delete from commentaires where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) GetComment(uid int) *[]structs.Commentaire { // Récupérer les données du commentaire grâce aux structs
	comment := structs.Commentaire{}
	rows, err := m.DB.Query("SELECT c.id, c.content, c.date, c.user_id, c.hidden, u.username, u.avatar FROM commentaires c LEFT JOIN users u ON u.id = c.user_id WHERE c.post_id=?", uid)
	tab := []structs.Commentaire{}
	var date int

	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&comment.Id, &comment.Content, &date, &comment.User.Id, &comment.Hidden, &comment.User.Username, &comment.User.Avatar)
		checkErr(err)
		comment.Date = m.DateConversion((date))
		comment.Likes = (*m.GetCommentLike(comment.Id))
		tab = append(tab, comment)
	}

	return &tab
}
func (m MyDB) GetCommentById(uid int) *structs.Commentaire { // Récupérer les données du commentaire grâce aux structs
	comment := structs.Commentaire{}
	rows, err := m.DB.Query("SELECT c.id, c.user_id, c.post_id FROM commentaires c LEFT JOIN users u ON u.id = c.user_id WHERE c.id=?", uid)

	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&comment.Id, &comment.User.Id, &comment.Post.Id)
		checkErr(err)
	}

	return &comment
}
func (m MyDB) GetUserByComment(uid int) *structs.Commentaire { // Récupérer les données du commentaire grâce aux structs
	comment := structs.Commentaire{}
	rows, err := m.DB.Query("SELECT u.id FROM commentaires c LEFT JOIN users u ON u.id = c.user_id WHERE c.id=?", uid)
	tab := structs.Commentaire{}

	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&comment.User.Id)
		checkErr(err)
	}

	return &tab
}

//CATEGORIES=============================================================================================
func (m MyDB) CreateCategory(name string) bool { // Pouvoir créer une catégorie
	stmt, err := m.DB.Prepare("INSERT INTO categories(name) values(?)")
	checkErr(err)

	_, err = stmt.Exec(name)
	checkErr(err)

	return true
}
func (m MyDB) UpdateCategory(id int, name string) bool { // Pouvoir modifier le nom de la catégorie

	stmt, err := m.DB.Prepare("update categories set name=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(name, id)
	checkErr(err)

	return true
}
func (m MyDB) DeleteCategory(id int) bool { // Pouvoir supprimer une catégorie
	stmt, err := m.DB.Prepare("select id from commentaires where post_id in(select id posts from posts where categorie_id=?)")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	stmt, err = m.DB.Prepare("select id from posts where categorie_id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	stmt, err = m.DB.Prepare("delete from categories where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) GetCategory(id int) string { // Récupérer les données de la catégorie
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
func (m MyDB) GetCategoryId(name string) int { // Récupérer les données de la catégorie depuis son identifiant
	rows, err := m.DB.Query("SELECT id FROM categories where name=?", name)
	checkErr(err)
	var id int

	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&id)
		checkErr(err)
	}
	fmt.Printf("name => %s | id => %d", name, id)

	return id
}
func (m MyDB) GetAllCategory() *[]structs.Categorie { // Récupérer les données de toutes les catégories depuis les structs
	rows, err := m.DB.Query("SELECT id,name FROM categories")
	checkErr(err)
	var cats []structs.Categorie
	var cat structs.Categorie
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&cat.Id, &cat.Name)
		cats = append(cats, cat)
		checkErr(err)
	}
	return &cats
}

//===========================================================================
func (m MyDB) CreateUserCat(user_id int, cat_id int) bool { // Pouvoir ajouter un nouveau modo et lui assigner les catégories
	stmt, err := m.DB.Prepare("INSERT INTO userCat(user_id, cat_id) values(?, ?)")
	checkErr(err)

	_, err = stmt.Exec(user_id, cat_id)
	checkErr(err)

	return true
}
func (m MyDB) DeleteUserCat(user_id int, cat_id int) bool { // mise à jour des catégories du modo

	// stmt, err := m.DB.Prepare("")update userCat set cat_id=? where user_id=? and cat_id <> ?
	stmt, err := m.DB.Prepare("update userCat set cat_id=? where user_id=?")
	if err != nil {
		m.CreateUserCat(user_id, cat_id)
	}

	_, err = stmt.Exec(user_id, cat_id)
	fmt.Println("ERREUR DELETEUSERCAT !!!!!!!!!!")
	checkErr(err)

	return true

}

//TICKETS===========================================================================

func (m MyDB) CreateTicket(id int, content string, categorieId int) bool { // Pouvoir créer un ticket
	stmt, err := m.DB.Prepare("INSERT INTO tickets(content, user_id, categorie_id) values(?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(id, content, categorieId)
	checkErr(err)

	return true
}
func (m MyDB) OpenTicket(id int) bool { // Pouvoir actualiser l'état d'un ticket en ouvert

	stmt, err := m.DB.Prepare("update categories set etat=1 where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) CloseTicket(id int) bool { // Pouvoir actualiser l'état d'un ticket en fermé

	stmt, err := m.DB.Prepare("update categories set etat=2 where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) GetAllTickt(etat int) *[]structs.Ticket { // Récupérer les données de tous les tickets depuis les structs
	tik := structs.Ticket{}
	tab := []structs.Ticket{}
	rows, err := m.DB.Query("SELECT t.id, t.content, t.date, t.etat, t.categorie_id, t.openBy, u.id user FROM tickets t LEFT JOIN users u ON t.user_id=u.id WHERE etat=? ORDER BY date ASC", etat)
	checkErr(err)
	var openBy int
	var user int

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&tik.Id, &tik.Content, &tik.Date, &tik.Etat, &tik.Categorie, &openBy, &user)
		checkErr(err)
		tik.User = *(m.GetUser(user))
		tik.OpenBy = *(m.GetUser(openBy))
		// tik.Date = time.Date(tik.Date)

		tNow := time.Unix(int64(1623234148), int64(0))

		tUnix := tNow.Unix()

		timeT := time.Unix(tUnix, 0)

		temp := timeT.String()

		fmt.Println(temp)

		temp = (strings.Split(temp, " "))[0]
		tab3 := strings.Split(temp, "-")
		var tab2 []string
		tab2 = append(tab2, tab3[2])
		tab2 = append(tab2, tab3[1])
		tab2 = append(tab2, tab3[0])
		tik.Date = strings.Join(tab2, "/")

		tab = append(tab, tik)
	}

	return &tab
}
func (m MyDB) GetTicket(id int) *structs.Ticket { // Récupérer les données d'un unique ticket
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

//MOTS DE PASSE=============================================================================
func hashMdp(mdp string) (string, error) { // Pouvoir chiffrer le mot de passe
	bytes, err := bcrypt.GenerateFromPassword([]byte(mdp), hashCost)
	return string(bytes), err
}
func (m MyDB) CompareMdp(password string, mail string) (int, error) { // Pouvoir comparer les mots de passe (ancien/nouveau) quand on le modifie
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
func (m MyDB) updateMdp(old string, mdp string, mail string) bool { // Pouvoir actualiser son mdp
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
func checkErr(err error) { // Pouvoir lancer un panic quand une erreur survient (gestion d'erreur)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
func (m MyDB) UpdateMdpForgotten(id int, password string) bool { // Récupération du mdp oublié
	pwd, err := hashMdp(password)
	checkErr(err)

	fmt.Println("HASH:", pwd)

	stmt, err := m.DB.Prepare("update users set mdp=? where id=?")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(pwd, id)
	checkErr(err)

	return true
}
func (m MyDB) DateConversion(date int) string { // Convertir le temps Unix en heure lisible

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

	return temp
}

func (m MyDB) GetStats(nb int) (string, error) { // Récupérer les statistiques des catégories
	var result string
	if nb == 0 {
		result = "date,nb"
		query := "select date from posts order by date asc"
		rows, err := m.DB.Query(query)
		checkErr(err)
		var dateF []string
		var counts []int
		count := 0
		var dates []int
		var date int

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&date)
			checkErr(err)
			dates = append(dates, date)
		}

		start := dates[0] - dates[0]%86400
		for i := 0; i < len(dates); i++ {
			if dates[i] > start+86400 {
				counts = append(counts, count)
				dateF = append(dateF, m.intToDate(start))
				fmt.Printf("%s => %d\n ", m.intToDate(start), count)
				for start < dates[i] {
					start += 86400
				}
				count = 0
			}
			count++
		}
		// fmt.Println(counts)
		// fmt.Println(dateF)
		for i := 0; i < len(counts); i++ {
			result += "|" + dateF[i] + "," + strconv.Itoa(counts[i])
		}
		fmt.Println(result)
	} else {
		rows, err := m.DB.Query("SELECT id, name from categories")
		checkErr(err)
		var ids []int
		var names []string
		var id int
		var name string
		result = "Category,Users,Modos,Admins"

		defer rows.Close()
		for rows.Next() {
			err = rows.Scan(&id, &name)
			checkErr(err)
			names = append(names, name)
			ids = append(ids, id)
		}

		wait := sync.WaitGroup{}
		lock := sync.Mutex{}
		for i := 0; i < len(ids); i++ {
			idstat := ids[i]
			name := names[i]
			wait.Add(1)
			go func(index int) {
				rows, err := m.DB.Query("select Users, Modos, Admins from (select count(p.id) Admins from posts p inner join users u on p.user_id=u.id inner join roles r on u.role_id=r.id where r.name like \"admin\" and p.categorie_id=?), (select count(p.id) Modos from posts p inner join users u on p.user_id=u.id inner join roles r on u.role_id=r.id where r.name like \"modo\" and p.categorie_id=?), (select count(p.id) Users from posts p inner join users u on p.user_id=u.id inner join roles r on u.role_id=r.id where r.name like \"user\" and p.categorie_id=?);", idstat, idstat, idstat)
				checkErr(err)
				var user int
				var modo int
				var admin int

				defer rows.Close()

				lock.Lock()
				for rows.Next() {
					err = rows.Scan(&user, &modo, &admin)
					checkErr(err)

					result += "|" + name + "," + strconv.Itoa(user) + "," + strconv.Itoa(modo) + "," + strconv.Itoa(admin)
				}

				lock.Unlock()
				wait.Done()
			}(i)
		}
		wait.Wait()
		fmt.Println(result)
	}
	return result, nil
}

func (m MyDB) Search(txt string, cats []string, limit int, offset int) *[]structs.Post { // Pouvoir rechercher par mot clé dans tous le site et filtrer les catégories
	fmt.Printf("search for : %s\n", txt)
	offset = offset * limit
	in := ""
	if len(cats) == 1 {
		in = " and p.categorie_id in(select id from categories where name like'" + cats[0] + "') "
	} else if len(cats) > 0 {
		in = " and p.categorie_id in(select id from categories where name in("
		for i := 0; i < len(cats); i++ {

			if i != 0 {
				in += ","
			}
			in += "\"" + cats[i] + "\""
		}
		in += ")) "
	}
	if len(txt) > 0 {
		txt = "and p.content like '%" + txt + "%'"
	}
	query := "SELECT p.id, p.content, p.date, p.categorie_id, p.hidden, p.user_id, u.username, u.avatar FROM posts p LEFT JOIN users u ON u.id = p.user_id WHERE p.hidden!=1 " + txt + " " + in + "ORDER BY date desc LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)
	fmt.Println(query)
	rows, err := m.DB.Query(query)
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
		post.Likes = (*m.GetPostLike(post.Id))
		tab = append(tab, post)
	}

	fmt.Println(tab)

	return &tab
}

func (m MyDB) intToDate(unix int) string {
	// fmt.Println(unix)
	result := ""
	result = strings.Split(time.Unix(int64(unix), 0).String(), " ")[0]
	// fmt.Println(tm)
	return result
}
