package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/github"
)

// GetMyNotifications prints your own notifications
func GetMyNotifications() error {

	fmt.Println()

	opt := &github.NotificationListOptions{
		All:           true,
		Participating: true,
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	var allItems []*github.Notification
	for {
		list, resp, err := client.Activity.ListNotifications(ctx, opt)
		if err != nil {
			return err
		}
		allItems = append(allItems, list...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, n := range allItems {
		fmt.Printf("[%s] %s:%s - %s (%s)\n",
			*n.Repository.Name, *n.Subject.Type, *n.Reason, *n.Subject.Title, *n.URL)
	}
	fmt.Println()
	return nil

}

// PrintUser prints one user
func PrintUser(username string) error {

	if username == "" {
		return fmt.Errorf("user argument required")
	}

	fmt.Printf("\nGetting user: %s\n\n", username)
	usr, _, err := client.Users.Get(ctx, username)
	if err != nil {
		return err
	}

	fmt.Printf("ID: %d\n", usr.ID)
	fmt.Printf("Name: %s\n", usr.GetName())
	fmt.Printf("Login: %s\n", usr.GetLogin())
	fmt.Printf("Email: %s\n", usr.GetEmail())
	fmt.Printf("Location: %s\n", usr.GetLocation())
	fmt.Printf("Created: %v\n", usr.GetCreatedAt())
	fmt.Printf("Company: %s\n", usr.GetCompany())

	fmt.Println()

	return nil
}

// PrintTeams prints teams and its members
func PrintTeams(org string) error {

	if org == "" {
		return fmt.Errorf("org argument required")
	}

	fmt.Println()

	opt := &github.ListOptions{PerPage: 10}
	var allItems []*github.Team
	for {
		list, resp, err := client.Organizations.ListTeams(ctx, org, opt)
		if err != nil {
			return err
		}
		allItems = append(allItems, list...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for _, e := range allItems {
		fmt.Printf("%d - %v\n", e.GetID(), *e.Name)
	}
	fmt.Println()
	return nil
}

// AddUserToTeam adds user to the specified team
func AddUserToTeam(teamID int64, username string) error {

	// validation
	if teamID == 0 || username == "" {
		log.Fatal("required argument missing")
	}
	// end validation

	// username
	usr, _, err := client.Users.Get(ctx, username)
	if err != nil {
		return err
	}
	// end user

	// team
	team, _, err := client.Organizations.GetTeam(ctx, teamID)
	if err != nil {
		return err
	}
	// end team

	// prompt
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Are you sure you want to add '%s' (%s) to '%s' team?: [Y/n]",
		username, usr.GetName(), team.GetName())
	resp, _ := reader.ReadString('\n')
	if resp != "Y\n" {
		return nil
	}
	//end prompt

	// is already member
	isMember, _, err := client.Organizations.IsTeamMember(ctx, teamID, username)
	if err != nil {
		return err
	}
	if isMember {
		fmt.Printf("%s already member of this team", username)
		return nil
	}
	// end if member

	// add user
	opt := &github.OrganizationAddTeamMembershipOptions{}
	_, _, err = client.Organizations.AddTeamMembership(ctx, teamID, username, opt)
	if err != nil {
		return err
	}
	fmt.Printf("%s has been added to this team", username)
	// end add user

	return nil

}
