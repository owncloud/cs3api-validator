Feature: Etag propagation
  As as a oCIS user
  I want to be edit files in my directory tree
  And let the clients know by changing the root etag

  Background: Public share exists
    Given user "admin" has logged in with password "admin"
    Given user "admin" has created a folder "a-folder" in the home directory with the alias "a-folder"
    Given user "admin" has created a folder "a-folder/a-sub-folder" in the home directory with the alias "a-sub-folder"

  Scenario: Change etag of personal home
    When user "admin" remembers the etag of the space with name "Admin"
    And user "admin" has uploaded a file "a-folder/a-sub-folder/testfile.txt" with content "text" in the home directory with the alias "testfile.txt"

