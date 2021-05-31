package structs

import "time"

type Session struct {
	Uuid			int
	User_Id			int

}

type Cookie struct {
	Name       		string
	Value      		string
	Path       		string
	Domain     		string
	Expires    		time.Time
	RawExpires 		string
// MaxAge=0 means no 'Max-Age' attribute specified.
// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
// MaxAge>0 means Max-Age attribute present and given in seconds
	MaxAge   		int
	Secure   		bool
	HttpOnly 		bool
	Raw      		string
	Unparsed 		[]string // Raw text of unparsed attribute-value pairs
}

type User struct {
	Id				int
	Username		string
	Mail			string
	Avatar			string
	SessionToken	int
	Role			int
}

type Commentaire struct {
	Id				int
	Content			string
	Date			int
	UserId			User
	PostId			Post
	CommentId		int
}

type Categorie struct {
	Id				int
	Name			string
}

type Post struct {
	Id				int
	Content 		string
	Date 			int
	UserId 			int
	CategorieId 	int
	Hidden 			bool
}

type Autorisation struct {
	Id				int
	Name 			string
}

type RoleAuth struct {
	AutorisationId Autorisation
	RoleId			Role
}

type CommentLike struct {
	UserId			User
	CommentaireId	Commentaire
	Vote			int
}

type PostLike struct {
	UserId			User
	PostId			Post
	Vote			int
}

type Badges struct {
	Id				int
	Name			string
	Image 			string
}

type BadgeUser struct {
	UserId			User
	BadgeId			Badges
}

type Role struct {
	Id				int
	Name			string
}

type Ticket struct {
	Id 				int
	Content			string
	Date			string
	Etat			bool
	UserId			User	
}

type BanList struct {
	StartDate		string
	EndDate			string
	Raison			string
	BanDef			string
	BannedBy		User
	UserId			User	
}

type All struct {
	Users	User
	Posts	Post
}