Feature: Etag and Treesize propagation
  As a user
  I want to edit files in my directory tree
  And let the clients know by changing the root etag

  Background: Personal space and some folders exist
    Given user "admin" has logged in with password "admin"
    Given user "admin" has created a personal space with the alias "Admin Home"
    Given user "admin" has created a folder "a-folder" in the home directory with the alias "a-folder"
    Given user "admin" has created a folder "a-folder/a-sub-folder" in the home directory with the alias "a-sub-folder"

  Scenario: Change etag of personal home and for the full tree
    Given user "admin" remembers the fileinfo of the resource with the alias "Admin Home"
    And user "admin" remembers the fileinfo of the resource with the alias "a-folder"
    And user "admin" remembers the fileinfo of the resource with the alias "a-sub-folder"
    When user "admin" has uploaded a file "a-folder/a-sub-folder/testfile.txt" with content "text" in the home directory with the alias "testfile.txt"
    Then for user "admin" the etag of the resource with the alias "Admin Home" should have changed
    And for user "admin" the etag of the resource with the alias "a-folder" should have changed
    And for user "admin" the etag of the resource with the alias "a-sub-folder" should have changed

  Scenario: Change etag of personal home and for only one subtree
    Given user "admin" has created a folder "a-folder/a-sub-folder-B" in the home directory with the alias "a-sub-folder-B"
    And user "admin" remembers the fileinfo of the resource with the alias "Admin Home"
    And user "admin" remembers the fileinfo of the resource with the alias "a-folder"
    And user "admin" remembers the fileinfo of the resource with the alias "a-sub-folder"
    And user "admin" remembers the fileinfo of the resource with the alias "a-sub-folder-B"
    When user "admin" has uploaded a file "a-folder/a-sub-folder/testfile.txt" with content "text" in the home directory with the alias "testfile.txt"
    Then for user "admin" the etag of the resource with the alias "Admin Home" should have changed
    And for user "admin" the etag of the resource with the alias "a-folder" should have changed
    And for user "admin" the etag of the resource with the alias "a-sub-folder" should have changed
    And for user "admin" the etag of the resource with the alias "a-sub-folder-B" should not have changed

  Scenario: Change treesize of personal home and for the full tree
    When user "admin" has uploaded a file "a-folder/a-sub-folder/testfile.txt" with content "text" in the home directory with the alias "testfile.txt"
    Then for user "admin" the treesize of the resource with the alias "Admin Home" should be 4
    When user "admin" has uploaded a file "a-folder/a-sub-folder/testfile2.txt" with content "text" in the home directory with the alias "testfile2.txt"
    And for user "admin" the treesize of the resource with the alias "Admin Home" should be 8
    And for user "admin" the treesize of the resource with the alias "a-folder" should be 8
    And for user "admin" the treesize of the resource with the alias "a-sub-folder" should be 8

  Scenario: Change treesize of personal home and for only one subtree
    Given user "admin" has created a folder "a-folder/a-sub-folder-B" in the home directory with the alias "a-sub-folder-B"
    When user "admin" has uploaded a file "a-folder/a-sub-folder/testfile.txt" with content "text" in the home directory with the alias "testfile.txt"
    Then for user "admin" the treesize of the resource with the alias "Admin Home" should be 4
    When user "admin" has uploaded a file "a-folder/a-sub-folder-B/testfile2.txt" with content "text" in the home directory with the alias "testfile2.txt"
    And for user "admin" the treesize of the resource with the alias "Admin Home" should be 8
    And for user "admin" the treesize of the resource with the alias "a-folder" should be 8
    And for user "admin" the treesize of the resource with the alias "a-sub-folder" should be 4
    And for user "admin" the treesize of the resource with the alias "a-sub-folder-B" should be 4
