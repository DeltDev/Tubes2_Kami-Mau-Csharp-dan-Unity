package BFS

import(
	"fmt"
	"backend/scraper"
)
func BFS(startPage string, endPage string) []string {

	path := [][]string{{startPage}}
	queue := []string{startPage}
	visited := make(map[string]bool)
	visited[startPage] = false
	if startPage == endPage {
		fmt.Println("Found the end page!")
		fmt.Println("Path: ", startPage)
		return []string{startPage}
	}
	var tempqueue []string
	var wannaGetLinks []string

	fmt.Println("flag 1")
	for len(queue) >= 0 {
		fmt.Println("flag 2")
		if len(queue) == 0 {

			if len(tempqueue) == 0 {
				fmt.Println("queue dan tempqueue habis")
				return []string{}
			}

			fmt.Println("queue habis")
			queue = append(queue, tempqueue...)
			tempqueue = []string{}
			fmt.Println("panjang queue: ", len(queue))
			fmt.Println("panjang tempqueue reset: ", len(tempqueue))
		}

		if len(queue) > 0 {
			if len(queue) > 4000 {
				wannaGetLinks = queue[:4000]
			} else {
				wannaGetLinks = queue
			}
			queue = queue[len(wannaGetLinks):]

			// fmt.Println("getLink jalan=============================================")
			parentAndChildMap := scraper.GetLinksMap(wannaGetLinks)
			fmt.Println("getLink selesai")
			fmt.Println("panjang parentAndChildMap: ", len(parentAndChildMap))
			fmt.Println("panjang queue: ", len(queue))
			// fmt.Println("parentAndChildMap: Asli ****************************** ", parentAndChildMap)
			parentAndChildMap = scraper.RemoveRedundantMap(parentAndChildMap)
			parentAndChildMap = scraper.RemoveRedundantMaphashtag(parentAndChildMap)

			for parent, arrChild := range parentAndChildMap {
				//fmt.Println("parent: ", parent)
				if parent == "Indonesia" {
					fmt.Println("================================")
				}
				for _, l := range arrChild {
					if !visited[l] {
						visited[l] = true
						tempqueue = append(tempqueue, l)
						foundParent := false
						i := 0
						for !foundParent {
							if path[i][len(path[i])-1] == parent {
								foundParent = true
								newPath := make([]string, len(path[i]))
								copy(newPath, path[i])
								newPath = append(newPath, l)
								path = append(path, newPath)
								if l == endPage {
									fmt.Println("Found the end page!")
									fmt.Println("Path: ", newPath)
									return newPath
								}
							}
							i++
						}

					}

				}
			}
		}
	}
	return []string{}
}