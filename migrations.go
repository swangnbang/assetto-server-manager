package servermanager

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"html/template"
)

const (
	CurrentMigrationVersion = 6
	versionMetaKey          = "version"
)

func Migrate(store Store) error {
	var storeVersion int

	err := store.GetMeta(versionMetaKey, &storeVersion)

	if err != nil && err != ErrMetaValueNotSet {
		return err
	}

	for i := storeVersion; i < CurrentMigrationVersion; i++ {
		err := migrations[i](store)

		if err != nil {
			return err
		}
	}

	return store.SetMeta(versionMetaKey, CurrentMigrationVersion)
}

type migrationFunc func(Store) error

var migrations = []migrationFunc{
	addEntrantIDToChampionships,
	addAdminAccount,
	championshipLinksToSummerNote,
	addEntrantsToChampionshipEvents,
	addIDToChampionshipClasses,
	enhanceOldChampionshipResultFiles,
}

func addEntrantIDToChampionships(rs Store) error {
	logrus.Infof("Running migration: Add Internal UUID to Championship Entrants")

	championships, err := rs.ListChampionships()

	if err != nil {
		return err
	}

	for _, c := range championships {
		for _, class := range c.Classes {
			for _, entrant := range class.Entrants {
				if entrant.InternalUUID == uuid.Nil {
					entrant.InternalUUID = uuid.New()
				}
			}
		}

		err = rs.UpsertChampionship(c)

		if err != nil {
			return err
		}
	}

	return nil
}

func addAdminAccount(rs Store) error {
	logrus.Infof("Running migration: Add Admin Account")

	account := NewAccount()
	account.Name = adminUserName
	account.DefaultPassword = "servermanager"
	account.Group = GroupAdmin

	return rs.UpsertAccount(account)
}

func championshipLinksToSummerNote(rs Store) error {
	logrus.Infof("Converting old championship links to new markdown format")

	championships, err := rs.ListChampionships()

	if err != nil {
		return err
	}

	var i int

	for _, c := range championships {
		i = 0
		for link, name := range c.Links {
			c.Info += template.HTML("<a href='" + link + "'>" + name + "</a>")

			if i != len(c.Links)-1 {
				c.Info += ", "
			}

			i++
		}

		c.Links = nil

		err = rs.UpsertChampionship(c)

		if err != nil {
			return err
		}
	}

	return nil
}

func addEntrantsToChampionshipEvents(rs Store) error {
	logrus.Infof("Running migration: Add entrants to championship events")

	championships, err := rs.ListChampionships()

	if err != nil {
		return err
	}

	for _, c := range championships {
		for _, event := range c.Events {
			if event.EntryList == nil || len(event.EntryList) == 0 {
				event.EntryList = c.AllEntrants()
			}
		}

		err = rs.UpsertChampionship(c)

		if err != nil {
			return err
		}
	}

	return nil
}

func addIDToChampionshipClasses(rs Store) error {
	logrus.Infof("Running migration: Add class ID to championship classes")

	championships, err := rs.ListChampionships()

	if err != nil {
		return err
	}

	for _, c := range championships {
		for _, class := range c.Classes {
			if class.ID == uuid.Nil {
				class.ID = uuid.New()
			}
		}

		err = rs.UpsertChampionship(c)

		if err != nil {
			return err
		}
	}

	return nil
}

func enhanceOldChampionshipResultFiles(rs Store) error {
	logrus.Infof("Running migration: Enhance Old Championship Results Files")

	championships, err := rs.ListChampionships()

	if err != nil {
		return err
	}

	for _, c := range championships {
		for _, event := range c.Events {
			for _, session := range event.Sessions {
				c.EnhanceResults(session.Results)
			}
		}

		err = rs.UpsertChampionship(c)

		if err != nil {
			return err
		}
	}

	return nil
}
