Feature: concurrent user sharing
  As a user
  I want to share their resource concurrently to different multiple users


  Scenario: users make concurrent sharing to each other
    Given user "admin" has logged in with password "admin"
    And user "admin" has created a personal space with the alias "Admin Home"
    And user "admin" has uploaded a file "testfile.txt" with content "concurrent sharing" in the home directory with the alias "testfile1"
    When user "admin" shares a file "testfile.txt" with the following users concurrently
      | users     |
      | marie     |
      | moss      |
      | richard   |
      | katherine |
      | einstein  |
    Then the concurrent user sharing should have been successfull
