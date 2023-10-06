Feature: concurrent user sharing
  As a users
  Users want to share their resource concurrently

  @concurrent
  Scenario: users make concurrent sharing to each other
    Given user "admin" has logged in with password "admin"
    And user "admin" has created a personal space with the alias "Admin Home"
    And user "admin" has uploaded a file "testfile1.txt" with content "concurrent sharing" in the home directory with the alias "testfile1"
    And user "admin" has uploaded a file "testfile2.txt" with content "concurrent sharing" in the home directory with the alias "testfile2"
    When user "admin" shares a file "testfile1.txt" with the following users concurrently
      | users     |
      | marie     |
      | moss      |
      | richard   |
      | katherine |
      | einstein  |
    Then the concurrent user sharing should have been successfull
