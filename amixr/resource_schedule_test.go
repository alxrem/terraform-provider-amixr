package amixr

import (
	"fmt"
	amixr "github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAmixrSchedule_basic(t *testing.T) {
	scheduleName := fmt.Sprintf("schedule-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAmixrScheduleResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAmixrScheduleConfig(scheduleName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAmixrScheduleResourceExists("amixr_schedule.test-acc-schedule"),
				),
			},
		},
	})
}

func testAccCheckAmixrScheduleResourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*amixr.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "amixr_schedule" {
			continue
		}

		if _, _, err := client.Schedules.GetSchedule(r.Primary.ID, &amixr.GetScheduleOptions{}); err == nil {
			return fmt.Errorf("Schedule still exists")
		}

	}
	return nil
}

func testAccAmixrScheduleConfig(scheduleName string) string {
	return fmt.Sprintf(`
resource "amixr_schedule" "test-acc-schedule" {
	name = "%s"
	time_zone = "America/New_York"
}
`, scheduleName)
}

func testAccCheckAmixrScheduleResourceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Schedule ID is set")
		}

		client := testAccProvider.Meta().(*amixr.Client)

		found, _, err := client.Schedules.GetSchedule(rs.Primary.ID, &amixr.GetScheduleOptions{})
		if err != nil {
			return err
		}
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Schedule policy not found: %v - %v", rs.Primary.ID, found)
		}
		return nil
	}
}
