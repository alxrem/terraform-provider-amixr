package amixr

import (
	"fmt"
	amixr "github.com/alertmixer/amixr-go-client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccAmixrOnCallShift_basic(t *testing.T) {
	scheduleName := fmt.Sprintf("schedule-%s", acctest.RandString(8))
	shiftName := fmt.Sprintf("shift-%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAmixrOnCallShiftResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAmixrOnCallShiftConfig(scheduleName, shiftName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAmixrOnCallShiftResourceExists("amixr_on_call_shift.test-acc-on_call_shift"),
				),
			},
		},
	})
}

func testAccCheckAmixrOnCallShiftResourceDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*amixr.Client)
	for _, r := range s.RootModule().Resources {
		if r.Type != "amixr_on_call_shift" {
			continue
		}

		if _, _, err := client.OnCallShifts.GetOnCallShift(r.Primary.ID, &amixr.GetOnCallShiftOptions{}); err == nil {
			return fmt.Errorf("OnCallShift still exists")
		}

	}
	return nil
}

func testAccAmixrOnCallShiftConfig(scheduleName string, shiftName string) string {
	return fmt.Sprintf(`
resource "amixr_schedule" "test-acc-schedule" {
	name = "%s"
}

resource "amixr_on_call_shift" "test-acc-on_call_shift" {
	name = "%s"
	schedule_id = amixr_schedule.test-acc-schedule.id
	type = "recurrent_event"
	start = "2020-09-04T16:00:00"
	duration = 3600
	level = 1
	frequency = "weekly"
	week_start = "SU"
	interval = 2
	by_day = ["MO", "FR"]
}
`, scheduleName, shiftName)
}

func testAccCheckAmixrOnCallShiftResourceExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No OnCallShift ID is set")
		}

		client := testAccProvider.Meta().(*amixr.Client)

		found, _, err := client.OnCallShifts.GetOnCallShift(rs.Primary.ID, &amixr.GetOnCallShiftOptions{})
		if err != nil {
			return err
		}
		if found.ID != rs.Primary.ID {
			return fmt.Errorf("OnCallShift policy not found: %v - %v", rs.Primary.ID, found)
		}
		return nil
	}
}
