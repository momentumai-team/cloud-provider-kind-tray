package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/getlantern/systray"
	"github.com/momentumai-team/cloud-provider-kind-tray/actions"
	"github.com/skratchdot/open-golang/open"
	v1 "k8s.io/api/core/v1"
)

var loadBalancerMenuItems = []*systray.MenuItem{}
var loadBalancersInUse = map[string]*systray.MenuItem{}

func main() {
	onExit := func() {
	}
	systray.Run(onReady, onExit)
}

func onReady() {
	var cmd *exec.Cmd
	var err error
	systray.SetIcon(getIcon("assets/load-balancer.png"))
	systray.SetTooltip("Cloud Provider Kind Loadbalancer")

	mStartOrig := systray.AddMenuItem("Start", "Start Kind Loadbalancer")
	systray.AddSeparator()
	loadBalancerMenu := systray.AddMenuItem("LoadBalancers", "LoadBalancers")
	loadBalancerMenu.Hide()
	for i := 0; i < 20; i++ {
		loadBalancerMenuItems = append(loadBalancerMenuItems, loadBalancerMenu.AddSubMenuItem("", ""))
		loadBalancerMenuItems[i].Hide()
	}
	systray.AddSeparator()
	mStopOrig := systray.AddMenuItem("Stop", "Stop Kind Loadbalancer")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	mStopOrig.Disable()

	stdoutChan := make(chan string)
	stderrChan := make(chan string)

	go func() {
		for {
			select {
			case <-mStartOrig.ClickedCh:
				loadBalancersInUse = map[string]*systray.MenuItem{}
				for _, item := range loadBalancerMenuItems {
					item.SetTitle("")
					item.Hide()
				}
				var stdout, stderr io.ReadCloser
				mStartOrig.Disable()
				mStopOrig.Enable()
				loadBalancerMenu.Enable()

				cmd, stdout, stderr, err = actions.Start()
				if err != nil {
					log.Fatal(err)
				}
				//Stream stdout
				go func() {
					scanner := bufio.NewScanner(stdout)
					for scanner.Scan() {
						stdoutChan <- scanner.Text()
					}
					close(stdoutChan)
				}()

				//Stream stderr
				go func() {
					scanner := bufio.NewScanner(stderr)
					for scanner.Scan() {
						stderrChan <- scanner.Text()
					}
					close(stderrChan)
				}()
			case <-mStopOrig.ClickedCh:
				loadBalancerMenu.Disable()
				mStopOrig.Disable()
				err = actions.Stop(cmd)
				if err != nil {
					log.Fatal(err)
				}
				mStartOrig.Enable()
				cmd = nil
			case <-loadBalancerMenuItems[0].ClickedCh:
				openWebView(0)
			case <-loadBalancerMenuItems[1].ClickedCh:
				openWebView(1)
			case <-loadBalancerMenuItems[2].ClickedCh:
				openWebView(2)
			case <-loadBalancerMenuItems[3].ClickedCh:
				openWebView(3)
			case <-loadBalancerMenuItems[4].ClickedCh:
				openWebView(4)
			case <-loadBalancerMenuItems[5].ClickedCh:
				openWebView(5)
			case <-loadBalancerMenuItems[6].ClickedCh:
				openWebView(6)
			case <-loadBalancerMenuItems[7].ClickedCh:
				openWebView(7)
			case <-loadBalancerMenuItems[8].ClickedCh:
				openWebView(8)
			case <-loadBalancerMenuItems[9].ClickedCh:
				openWebView(9)
			case <-loadBalancerMenuItems[10].ClickedCh:
				openWebView(10)
			case <-loadBalancerMenuItems[11].ClickedCh:
				openWebView(11)
			case <-loadBalancerMenuItems[12].ClickedCh:
				openWebView(12)
			case <-loadBalancerMenuItems[13].ClickedCh:
				openWebView(13)
			case <-loadBalancerMenuItems[14].ClickedCh:
				openWebView(14)
			case <-loadBalancerMenuItems[15].ClickedCh:
				openWebView(15)
			case <-loadBalancerMenuItems[16].ClickedCh:
				openWebView(16)
			case <-loadBalancerMenuItems[17].ClickedCh:
				openWebView(17)
			case <-loadBalancerMenuItems[18].ClickedCh:
				openWebView(18)
			case <-loadBalancerMenuItems[19].ClickedCh:
				openWebView(19)
			case <-mQuitOrig.ClickedCh:
				mStopOrig.Disable()
				mQuitOrig.Disable()
				mStartOrig.Disable()
				if cmd != nil {
					err = actions.Stop(cmd)
					if err != nil {
						log.Fatal(err)
					}
				}
				systray.Quit()
				return
			}

		}
	}()
	// Read from channels and print to stdout
	for stdoutChan != nil || stderrChan != nil {
		select {
		case out, ok := <-stdoutChan:
			if !ok {
				stdoutChan = make(chan string)
				continue
			}
			lbEventStatus := make(map[string]v1.LoadBalancerStatus)
			err := json.Unmarshal([]byte(out), &lbEventStatus)
			if err == nil {
				for eventStatus, status := range lbEventStatus {
					for _, ingress := range status.Ingress {
						for _, lb := range ingress.Ports {
							key := fmt.Sprintf("%s:%d", ingress.IP, lb.Port)
							if strings.EqualFold(eventStatus, "added") || strings.EqualFold(eventStatus, "updated") {
								if ok, lbMenuItem := doesLBMenuExist(key); ok {
									lbMenuItem.Show()
								} else {
									lbMenuItem := findEmptyMenuItem()
									if lbMenuItem != nil {
										loadBalancerMenu.Show()
										lbMenuItem.SetTitle(key)
										lbMenuItem.Show()
										loadBalancersInUse[key] = lbMenuItem
									}
								}
							} else if strings.EqualFold(eventStatus, "deleted") {
								if lbMenuItem, ok := loadBalancersInUse[key]; ok {
									lbMenuItem.Hide()
								}
							}
						}
					}
				}
			} else {
				fmt.Println("STDOUT:", out)
			}
		case _, ok := <-stderrChan:
			if !ok {
				stderrChan = make(chan string)
				continue
			}
			//fmt.Println("STDERR:", err)
		}
	}
}

func openWebView(index int) {
	menuItem := loadBalancerMenuItems[index]
	url := findMenuItemUrl(menuItem)
	err := open.Run(fmt.Sprintf("http://%s", url))
	if err != nil {
		fmt.Println(err)
	}
}

func doesLBMenuExist(key string) (bool, *systray.MenuItem) {
	if lbMenuItem, ok := loadBalancersInUse[key]; ok {
		return true, lbMenuItem
	} else {
		return false, nil
	}
}

func findMenuItemUrl(menuItem *systray.MenuItem) string {
	for key, lbMenuItem := range loadBalancersInUse {
		if lbMenuItem == menuItem {
			return key
		}
	}
	return ""
}

func findEmptyMenuItem() *systray.MenuItem {
	for _, lbMenuItem := range loadBalancerMenuItems {
		if !doesItemExistInMap(lbMenuItem) {
			return lbMenuItem
		}
	}
	return nil
}

func doesItemExistInMap(menuItem *systray.MenuItem) bool {
	for _, lbMenuItem := range loadBalancersInUse {
		if lbMenuItem == menuItem {
			return true
		}
	}
	return false
}

func getIcon(s string) []byte {
	b, err := os.ReadFile(s)
	if err != nil {
		fmt.Print(err)
	}
	return b
}
