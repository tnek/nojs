package admin

import (
	"fmt"

	"github.com/tnek/notes-site/model"
)

const (
	kAdminPassword = "BF9tGtjYi3ZRvwhPVg9Q"
)

var (
	firefoxNotes = map[string]string{
		"Welcome to your new digital bulletin board": "Share whatever you want on it or place a bulletin onto any other user's board. Just remember a few rules:<ul><li>No javascript allowed whatsoever!</li><li>No straws allowed</li></ul>",
		"Public Menace Alert":                        "<script>alert('weewoo')</script>I saw Kent at a grocery store in Los Angeles yesterday. I told him how cool it was to meet him in person, but I didn’t want to be a douche and bother him and ask him for photos or anything. He said, “Oh, like you’re doing now?” I was taken aback, and all I could say was “Huh?” but he kept cutting me off and going “huh? huh? huh?” and closing his hand shut in front of my face. I walked away and continued with my shopping, and I heard him chuckle as I walked off. When I came to pay for my stuff up front I saw him trying to walk out the doors with like fifteen Milky Ways in his hands without paying.<br>The girl at the counter was very nice about it and professional, and was like “Sir, you need to pay for those first on the Square Register™.” At first he kept pretending to be tired and not hear her, but eventually turned back around and brought them to the Square Register™ on the counter.<br>When she took one of the bars and started scanning it multiple times, he stopped her and told her to scan them each individually “to prevent any electrical infetterence,” and then turned around and winked at me. I don’t even think that’s a word. After she scanned each bar and put them in a bag and started to say the price, he kept interrupting her by yawning really loudly.<br>",

		"Hello!": "Hello it's me your admin!",
	}
)

func firefoxAdminName(u *model.User) string {
	return fmt.Sprintf("admin-%v", u.ID)
}

func initAdminUser(adminName string) (*model.User, error) {
	uid, err := model.NewUser(adminName, kAdminPassword, true)
	if err != nil {
		return nil, err
	}
	admin, err := model.UserByUUID(uid)
	if err != nil {
		return nil, err
	}

	return admin, nil
}

func (a *Admin) sharedNote(from *model.User, to *model.User, title string, contents string) error {
	note, err := model.NewNote(from, title, contents)
	if err != nil {
		return err
	}
	if _, err := model.ShareNote(from, note, to.Name); err != nil {
		return err
	}
	return nil
}

func (a *Admin) NewAdmin(u *model.User) ([]*model.User, error) {
	fn := firefoxAdminName(u)

	fu, err := initAdminUser(fn)
	if err != nil {
		return nil, err
	}

	for title, contents := range firefoxNotes {
		if err := a.sharedNote(fu, u, title, contents); err != nil {
			return nil, err
		}
	}

	if _, err := model.NewNote(fu, a.FirefoxFlag, a.FirefoxFlag); err != nil {
		return nil, err
	}

	return []*model.User{fu}, nil
}