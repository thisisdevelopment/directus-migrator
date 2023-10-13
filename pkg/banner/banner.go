package banner

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

func Display() {
	fmt.Printf(
		`------------------------------------------------------------------
Directus Migrator. This is a %s Project. https://this.nl
------------------------------------------------------------------
`, aurora.Yellow("TISGroup"))
}
