package structs

import "time"

type Session struct {
	Uuid    string
	User_Id int
}

type SameSite int

const (
	SameSiteDefaultMode SameSite = iota + 1
	SameSiteLaxMode
	SameSiteStrictMode
	SameSiteNoneMode
)

type Cookie struct {
	Name       string
	Value      string
	Path       string
	Domain     string
	Expires    time.Time
	RawExpires string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   int
	Secure   bool
	SameSite SameSite
	HttpOnly bool
	Raw      string
	Unparsed []string // Raw text of unparsed attribute-value pairs
}

type User struct {
	Id           int
	Username     string
	Mail         string
	Avatar       string
	SessionToken string
	Role         int
	Verif        int
}
type Users struct {
	Users []User
	Error bool
}
type UserToken struct {
	Userid int
	Token  string
}

type Commentaire struct {
	Id        int
	Content   string
	Date      string
	Hidden    bool
	User      User
	Post      Post
	CommentId int
	Likes     Like
}
type Commentaires struct {
	Commentaires []Commentaire
	Post         Post
	Error        bool
	User         User
	Page         string
}

type Categorie struct {
	Id   int
	Name string
}
type Categories struct {
	Categories []Categorie
	Error      bool
	User       User
	Page       string
}

type Post struct {
	Id        int
	Content   string
	Date      string
	User      User
	Categorie string
	Hidden    bool
	Likes     Like
}
type Posts struct {
	Posts []Post
	Error bool
	User  User
	Page  string
}
type Like struct {
	Note  int
	Total int
}

type Autorisation struct {
	Id   int
	Name string
}
type Autorisations struct {
	Autorisations []Autorisation
	Error         bool
}

type RoleAuth struct {
	AutorisationId Autorisation
	RoleId         Role
}
type RoleAuths struct {
	RoleAuths []RoleAuth
	Error     bool
}

type CommentLike struct {
	UserId        User
	CommentaireId Commentaire
	Vote          int
}
type Commentlikes struct {
	Commentlikes []CommentLike
	Error        bool
}

type PostLike struct {
	UserId User
	PostId Post
	Vote   int
}
type Postlikes struct {
	Postlikes []PostLike
	Error     bool
	Page      string
}

type Badge struct {
	Id    int
	Name  string
	Image string
}

type BadgeUser struct {
	User   User
	Badges []Badge
	Error  bool
	Page   string
}

type Role struct {
	Id   int
	Name string
}
type Roles struct {
	Error bool
	Modo  []User
	Users []User
	Admin []User
	User  User
	Page  string
}

type Ticket struct {
	Id        int
	Content   string
	Date      string
	Etat      int
	Categorie int
	User      User
	OpenBy    User
}
type Tickets struct {
	Wait  []Ticket
	Open  []Ticket
	Close []Ticket
	Error bool
	User  User
	Page  string
}

type BanList struct {
	Id        int
	StartDate string
	EndDate   string
	Raison    string
	BanDef    string
	BannedBy  User
	User      User
}
type BanLists struct {
	BanLists []BanList
	Error    bool
	User     User
	Page     string
}

type Err0r struct {
	Error bool
	User  User
	Page  string
}

type UserCat struct {
	Username   string
	Id         int
	Categories []int
}

type Stats struct {
	Seven string
	Month string
	All   string
}
