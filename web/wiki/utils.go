package wiki

import (
	"diskhub/web/models"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func OrganizeArticles(articles *[]models.DocElement, object models.FileObject) error {
	orgFile, err := object.Found(".org")
	if err != nil && err != models.ErrFileNotFound {
		return err
	}
	if orgFile.IsDir {
		return fmt.Errorf(".org need to be a file in %s", object.Name)
	}

	content, err := os.ReadFile(orgFile.Path)
	if err != nil && err != os.ErrNotExist {
		return err
	} else if err == os.ErrNotExist {
		return nil
	}

	lines := strings.Split(string(content), "\n")
	organization := make(map[string]int)
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			continue
		}

		name := strings.TrimSpace(parts[1])
		value, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return err
		}

		organization[name] = value
	}

	ordered := make([]models.DocElement, 0)
	unordered := make([]models.DocElement, 0)

	for _, art := range *articles {
		if _, ok := organization[art.GetName()]; ok {
			ordered = append(ordered, art)
		} else {
			unordered = append(unordered, art)
		}
	}

	// Trier les articles ordonnés selon organization
	sort.Slice(ordered, func(i, j int) bool {
		return organization[ordered[i].GetName()] < organization[ordered[j].GetName()]
	})

	// Combiner
	*articles = append(ordered, unordered...)

	return nil
}

func organizeGroup(group *models.DocGroup, object models.FileObject) error {
	// Convertir en []DocElement
	elements := make([]models.DocElement, len(group.Articles))
	for i := range group.Articles {
		elements[i] = &group.Articles[i] // ou group.Articles[i] selon votre type
	}

	// Appeler la fonction
	err := OrganizeArticles(&elements, object)
	if err != nil {
		return err
	}

	// Reconvertir
	group.Articles = []models.DocArticle{}
	for _, elem := range elements {
		group.Articles = append(group.Articles, *elem.(*models.DocArticle))
	}

	return nil
}
