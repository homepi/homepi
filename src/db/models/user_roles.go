package models

func (user *User) CanRunAccessory() bool {
	if user.Role.Administrator {
		return true
	}
	return user.Role.CanRunAccessory
}

func (user *User) CanSeeAccessories() bool {
	if user.Role.Administrator {
		return true
	}
	return user.Role.CanSeeAccessories
}

func (user *User) CanCreateAccessory() bool {
	if user.Role.Administrator {
		return true
	}
	return user.Role.CanCreateAccessory
}

func (user *User) CanRemoveAccessory() bool {
	if user.Role.Administrator {
		return true
	}
	return user.Role.CanRemoveAccessory
}

func (user *User) CanSeeWebhook() bool {
	if user.Role.Administrator {
		return true
	}
	return user.Role.CanSeeWebhook
}

func (user *User) CanCreateWebhook() bool {
	if user.Role.Administrator {
		return true
	}
	return user.Role.CanCreateWebhook
}

func (user *User) CanRemoveWebhook() bool {
	if user.Role.Administrator {
		return true
	}
	return user.Role.CanRemoveWebhook
}

func (user *User) CanSeeUsers() bool {
	if user.Role.Administrator {
		return true
	}
	return user.Role.CanSeeUsers
}

func (user *User) CanCreateUser() bool {
	if user.Role.Administrator {
		return true
	}
	return user.Role.CanSeeUsers
}

func (user *User) CanRemoveUser() bool {
	if user.Role.Administrator {
		return true
	}
	return user.Role.CanSeeUsers
}

func (user *User) CanSeeRoles() bool {
	if user.Role.Administrator {
		return true
	}
	return user.Role.CanSeeUsers
}

func (user *User) CanCreateRole() bool {
	if user.Role.Administrator {
		return true
	}
	return user.Role.CanSeeUsers
}

func (user *User) CanRemoveRole() bool {
	if user.Role.Administrator {
		return true
	}
	return user.Role.CanSeeUsers
}

func (user *User) CanSeeLogs() bool {
	if user.Role.Administrator {
		return true
	}
	return user.Role.CanSeeUsers
}
