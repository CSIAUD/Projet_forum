package bdd

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	// "strconv"
	structs "Forum/static/go/structs"
	"golang.org/x/crypto/bcrypt"
)

type MyDB struct {
	DB *sql.DB
}

const userRole = 3
const hashCost = 18

//===================================================================================================
func (m MyDB) CreateCommentLike(user_id int, commentid int, vote int) bool {
	stmt, err := m.DB.Prepare("INSERT INTO comments(user_id, comment_id, vote) values(?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(user_id, commentid, vote)
	checkErr(err)

	return true
}
func (m MyDB) UpdateCommentLike(id int, user_id int, comment_id int, vote int) bool {

	stmt, err := m.DB.Prepare("update comments set user_id=? comment_id=? vote=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(user_id, comment_id, vote, id)
	checkErr(err)

	return true
}
func (m MyDB) DeleteCommentLike(id int) bool {
	stmt, err := m.DB.Prepare("delete from comments where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) GetCommentLike(id int) {
	rows, err := m.DB.Query("SELECT user_id, comment_id, vote FROM users where id=?", id)
	checkErr(err)
	var user_id int
	var comment_id int
	var vote int

	if rows.Next() {
		err = rows.Scan(&user_id, &comment_id, &vote)
		checkErr(err)
		fmt.Println(user_id)
		fmt.Println(comment_id)
		fmt.Println(vote)
	}
}

func (m MyDB) CreatePostLike(user_id int, post_id int, vote int) bool {
	stmt, err := m.DB.Prepare("INSERT INTO posts(user_id, post_id, vote) values(?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(user_id, post_id, vote)
	checkErr(err)

	return true
}
func (m MyDB) UpdatePostLike(id int, user_id int, post_id int, vote int) bool {

	stmt, err := m.DB.Prepare("update posts set user_id=? post_id=? vote=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(user_id, post_id, vote, id)
	checkErr(err)

	return true
}
func (m MyDB) DeletePostLike(id int) bool {
	stmt, err := m.DB.Prepare("delete from posts where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) GetPostLike(id int) {
	rows, err := m.DB.Query("SELECT user_id, post_id, vote FROM users where id=?", id)
	checkErr(err)
	var user_id int
	var post_id int
	var vote int

	if rows.Next() {
		err = rows.Scan(&user_id, &post_id, &vote)
		checkErr(err)
		fmt.Println(user_id)
		fmt.Println(post_id)
		fmt.Println(vote)
	}
}

//===================================================================================================
func (m MyDB) AddBadgeUser(user_id int, badge_id int) bool {
	stmt, err := m.DB.Prepare("INSERT INTO badgeUser(user_id, badge_id) values(?,?)")
	checkErr(err)

	_, err = stmt.Exec(user_id, badge_id)
	checkErr(err)

	return true
}
func (m MyDB) DeleteBadgeUser(id int) bool {
	stmt, err := m.DB.Prepare("delete from badgeUser where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) GetBadgeUser(id int) (int, int) {
	rows, err := m.DB.Query("SELECT user_id, badge_id FROM badgeUser where id=?", id)
	checkErr(err)
	var user_id int
	var badge_id int

	if rows.Next() {
		err = rows.Scan(&user_id, &badge_id)
		checkErr(err)
	}

	return user_id, badge_id
}

func (m MyDB) GetBadge(id int) (string, string) {
	rows, err := m.DB.Query("SELECT name, image FROM badges where id=?", id)
	checkErr(err)
	var name string
	var image string

	if rows.Next() {
		err = rows.Scan(&name, &image)
		checkErr(err)
	}

	return name, image
}

func (m MyDB) GetAuth(id int) string {
	rows, err := m.DB.Query("SELECT name FROM autorisations where id=?", id)
	checkErr(err)
	var name string

	if rows.Next() {
		err = rows.Scan(&name)
		checkErr(err)
	}

	return name
}

func (m MyDB) GetRole(id int) string {
	rows, err := m.DB.Query("SELECT name FROM roles where id=?", id)
	checkErr(err)
	var name string

	if rows.Next() {
		err = rows.Scan(&name)
		checkErr(err)
	}

	return name
}

//===================================================================================================
func (m MyDB) CreatebanList(startDate string, endDate string, raison string, banDef string, bannedBy int, user_id int) bool {
	raison, err := hashMdp(raison)

	stmt, err := m.DB.Prepare("INSERT INTO banList(startDate, endDate, raison, banDef, bannedBy, user_id) values(?,?,?,?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(startDate, endDate, raison, banDef, bannedBy, user_id)
	checkErr(err)

	return true
}
func (m MyDB) UpdatebanList(startDate string, endDate string, raison string, banDef string, bannedBy string, user_id int) bool {
	stmt, err := m.DB.Prepare("update banList set startDate=?, endDate=?, raison=?, banDef=?  where bannedBy=?, user_id=?")
	checkErr(err)

	_, err = stmt.Exec(startDate, endDate, raison, banDef, bannedBy, user_id)
	checkErr(err)

	return true
}
func (m MyDB) DeletebanList(user_id int) bool {
	stmt, err := m.DB.Prepare("delete from banList where user_id=?")
	checkErr(err)

	_, err = stmt.Exec(user_id)
	checkErr(err)

	return true
}
func (m MyDB) GetbanList(user_id int) {
	rows, err := m.DB.Query("SELECT startDate,endDate,raison,banDef, FROM banList where user_id=?", user_id)
	checkErr(err)
	var startDate string
	var endDate string
	var raison string
	var banDef string
	if rows.Next() {
		err = rows.Scan(&startDate, &endDate, &raison, &banDef)
		checkErr(err)
		fmt.Println(startDate)
		fmt.Println(endDate)
		fmt.Println(raison)
		fmt.Println(banDef)
	}
}

//===================================================================================================
func (m MyDB) CreateUser(username string, mail string, mdp string, avatar string, sessionToken string, role_id int) bool {
	mdp, err := hashMdp(mdp)

	stmt, err := m.DB.Prepare("INSERT INTO users(username, mail, mdp, avatar, sessionToken, role_id) values(?,?,?,?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(username, mail, mdp, avatar, sessionToken, role_id)
	checkErr(err)

	return true
}
func (m MyDB) UpdateUser(username string, mail string, mdp string, avatar string, id int) bool {
	stmt, err := m.DB.Prepare("update users set name=?, firstName=?, username=?, mail=?, avatar=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(username, mail, mdp, avatar, id)
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
	rows, err := m.DB.Query("SELECT username,mail,mdp,avatar FROM users where id=?", id)
	checkErr(err)
	var username string
	var mail string
	var mdp string
	var avatar string

	if rows.Next() {
		err = rows.Scan(&username, &mail, &mdp, &avatar)
		checkErr(err)
	}
	val := structs.User{Username: username, Mail: mail, Mdp: mdp, Avatar: avatar}
	return &val
}

//==================================================================================================================
func (m MyDB) CreatePost(id int, content string, date int, hidden bool, userID int, categorieID int) bool {

	stmt, err := m.DB.Prepare("INSERT INTO posts(id, content, date, hidden, userID, categorieID) values(?,?,?,?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(id, content, time.Now(), hidden, userID, categorieID)
	checkErr(err)

	return true
}
func (m MyDB) UpdatePost(id int, content string, categorieID int, hidden bool) bool {

	if !(hidden) {
		hidden = false
	}
	stmt, err := m.DB.Prepare("update posts set content=?, hidden=?, categorie_id=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(content, hidden, categorieID, id)
	checkErr(err)

	return true
}
func (m MyDB) DeletePost(id int) bool {
	stmt, err := m.DB.Prepare("delete from posts where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) GetPost(uid int) {
	rows, err := m.DB.Query("SELECT id, content, date, user_id, categorie_id FROM posts where id=?", uid)
	checkErr(err)
	var id int
	var content string
	var date int
	var user_id int
	var categorie_id int

	if rows.Next() {
		err = rows.Scan(&id, &content, &date, &user_id, &categorie_id)
		checkErr(err)
	}
}


func (m MyDB) GetNbPost(limit int, offset int) *[]structs.Post {
	offset = offset * limit
	post := structs.Post{}

	rows, err := m.DB.Query("SELECT id, content, date, user_id, categorie_id, hidden FROM posts order by date LIMIT ? OFFSET ?",limit,offset)
	checkErr(err)
	tab := []structs.Post{}

	for rows.Next() {
		err = rows.Scan(&post.Id, &post.Content, &post.Date, &post.UserId, &post.CategorieId, &post.Hidden)
		checkErr(err)
		tab = append(tab, post)
	}

	fmt.Println(tab)
	rows.Close()

	return &tab
}

func (m MyDB) CreateComment(id int, content string, date int, UserId int, postId int, commentId int) bool {
	stmt, err := m.DB.Prepare("INSERT INTO commentaires(id, content, date, UserId, postId, commentId) values(?,?,?,?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(id, content, time.Now(), UserId, postId, commentId)
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
	stmt, err := m.DB.Prepare("delete from commentaires where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) GetComment(uid int) {
	rows, err := m.DB.Query("SELECT id, content, date, user_id, post_id, comment_id FROM users where id=?", uid)
	checkErr(err)
	var id int
	var content string
	var date int
	var by int
	var PostId int
	var CommentId int

	if rows.Next() {
		err = rows.Scan(&id, &content, &date, &by, &PostId, &CommentId)
		checkErr(err)
		fmt.Println(id)
		fmt.Println(content)
		fmt.Println(date)
		fmt.Println(by)
		fmt.Println(PostId)
		fmt.Println(CommentId)
	}
}

//==============================================================================================
func (m MyDB) CreateCategory(txt string) bool {
	stmt, err := m.DB.Prepare("INSERT INTO categories(txt) values(?)")
	checkErr(err)

	_, err = stmt.Exec(txt)
	checkErr(err)

	return true
}
func (m MyDB) UpdateCategory(id int, txt string) bool {

	stmt, err := m.DB.Prepare("update categories set txt=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(txt, id)
	checkErr(err)

	return true
}
func (m MyDB) DeleteCategory(id int) bool {
	stmt, err := m.DB.Prepare("delete from categories where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) GetCategory(id int) {
	rows, err := m.DB.Query("SELECT txt FROM users where id=?", id)
	checkErr(err)
	var txt string

	if rows.Next() {
		err = rows.Scan(&txt)
		checkErr(err)
		fmt.Println(txt)
	}
}

//===========================================================================
func (m MyDB) CreateRoleAuth(autorisationId int, roleId int) bool {
	stmt, err := m.DB.Prepare("INSERT INTO roleAuth(autorisation_id, badge_id) values(?,?)")
	checkErr(err)

	_, err = stmt.Exec(autorisationId, roleId)
	checkErr(err)

	return true
}
func (m MyDB) DeleteRoleAuth(id int) bool {
	stmt, err := m.DB.Prepare("delete from roleAuth where id=?")
	checkErr(err)

	_, err = stmt.Exec(id)
	checkErr(err)

	return true
}
func (m MyDB) GetRoleAuth(id int) {
	rows, err := m.DB.Query("SELECT autorisation_id, role_id FROM users where id=?", id)
	checkErr(err)
	var authId int
	var roleId int

	if rows.Next() {
		err = rows.Scan(&authId, &roleId)
		checkErr(err)
		fmt.Println(authId)
		fmt.Println(roleId)
	}
}

//===========================================================================
func (m MyDB) CreateTickets(id int, content string, date int, etat bool, userId int) bool {
	stmt, err := m.DB.Prepare("INSERT INTO Tickets(id, content, date, etat, user_id) values(?,?,?,?,?)")
	checkErr(err)

	_, err = stmt.Exec(id, content, date, etat, userId)
	checkErr(err)

	return true
}
func (m MyDB) UpdateTickets(id int, content string, date int, etat bool) bool {

	stmt, err := m.DB.Prepare("update categories set content=? date=? etat=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(id, content, date, etat)
	checkErr(err)

	return true
}
func (m MyDB) GetTickets(uid int) (int, string, int, bool, int) {
	rows, err := m.DB.Query("SELECT id, content, date, etat, user_id FROM users where id=?", uid)
	checkErr(err)
	var id int
	var content string
	var date int
	var etat bool
	var userId int

	if rows.Next() {
		err = rows.Scan(&id, &content, &date, &etat, &userId)
		checkErr(err)
	}
	return id, content, date, etat, userId
}

//==============================================================================
func hashMdp(mdp string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(mdp), hashCost)
	return string(bytes), err
}
func (m MyDB) compareMdp(password string, id int) bool {
	rows, err := m.DB.Query("SELECT mdp FROM users where id=?", id)
	checkErr(err)
	var mdp string

	if rows.Next() {
		err = rows.Scan(&mdp)
		checkErr(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(mdp), []byte(password))
	return err == nil
}
func (m MyDB) updateMdp(old string, mdp string, id int) bool {
	if !m.compareMdp(old, id) {
		return false
	}
	stmt, err := m.DB.Prepare("update users set mdp=? where id=?")
	checkErr(err)

	_, err = stmt.Exec(mdp, id)
	checkErr(err)
	return true
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}