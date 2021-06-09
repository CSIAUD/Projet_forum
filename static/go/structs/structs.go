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

type Commentaire struct {
	Id        int
	Content   string
	Date      string
	User      User
	Post      Post
	CommentId int
}
type Commentaires struct {
	Commentaires []Commentaire
	Error        bool
	User         User
}

type Categorie struct {
	Id   int
	Name string
}
type Categories struct {
	Categories []Categorie
	Error      bool
}

type Post struct {
	Id        int
	Content   string
	Date      string
	User      User
	Categorie string
	Hidden    bool
	Likes     int
}
type Posts struct {
	Posts []Post
	Error bool
	User  User
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
}

type Role struct {
	Id   int
	Name string
}
type Roles struct {
	Roles []Role
	Error bool
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
	Wait 	[]Ticket
	Open 	[]Ticket
	Close	[]Ticket
	Error   bool
	User 	User
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
}

type Err0r struct {
	Error bool
	User  User
}
