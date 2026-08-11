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

  Scenario: Change etag of personal home and move subtree
    Given user "admin" has created a folder "a-folder/a-sub-folder-B" in the home directory with the alias "a-sub-folder-B"
    And user "admin" has created a folder "a-folder/a-sub-folder/testFolder" in the home directory with the alias "testFolder"
    And user "admin" remembers the fileinfo of the resource with the alias "Admin Home"
    And user "admin" has created a folder "a-folder/a-sub-folder/testFolder/test2" in the home directory with the alias "test2"
    # This is needed to make sure that the propagation has happened to the testing file tree before we move the folder
    And for user "admin" the etag of the resource with the alias "Admin Home" should have changed
    And user "admin" remembers the fileinfo of the resource with the alias "Admin Home"
    And user "admin" has uploaded a file "a-folder/a-sub-folder/testfile.txt" with content "text" in the home directory with the alias "testfile.txt"
    # This is needed to make sure that the propagation has happened to the testing file tree before we move the folder
    And for user "admin" the etag of the resource with the alias "Admin Home" should have changed
    # Now we have a clean state and can start the actual test
    And user "admin" remembers the fileinfo of the resource with the alias "Admin Home"
    And user "admin" remembers the fileinfo of the resource with the alias "a-folder"
    And user "admin" remembers the fileinfo of the resource with the alias "a-sub-folder"
    And user "admin" remembers the fileinfo of the resource with the alias "a-sub-folder-B"
    And user "admin" remembers the fileinfo of the resource with the alias "testFolder"
    And user "admin" remembers the fileinfo of the resource with the alias "test2"
    When user "admin" moves the resource with alias "a-sub-folder" inside a space to target "a-folder/a-sub-folder-B/a-sub-folder"
    Then for user "admin" the etag of the resource with the alias "Admin Home" should have changed
    And for user "admin" the etag of the resource with the alias "a-folder" should have changed
    And for user "admin" the etag of the resource with the alias "a-sub-folder" should not have changed
    And for user "admin" the etag of the resource with the alias "a-sub-folder-B" should have changed
    And for user "admin" the etag of the resource with the alias "testFolder" should not have changed
    And for user "admin" the etag of the resource with the alias "test2" should not have changed

  Scenario: Change treesize of personal home and for the full tree
    When user "admin" has uploaded a file "a-folder/a-sub-folder/testfile.txt" with content "text" in the home directory with the alias "testfile.txt"
    Then for user "admin" the treesize of the resource with the alias "Admin Home" should be 4
    When user "admin" has uploaded a file "a-folder/a-sub-folder/testfile2.txt" with content "text" in the home directory with the alias "testfile2.txt"
    Then for user "admin" the treesize of the resource with the alias "Admin Home" should be 8
    And for user "admin" the treesize of the resource with the alias "a-folder" should be 8
    And for user "admin" the treesize of the resource with the alias "a-sub-folder" should be 8

  Scenario: Change treesize of personal home and for only one subtree
    Given user "admin" has created a folder "a-folder/a-sub-folder-B" in the home directory with the alias "a-sub-folder-B"
    When user "admin" has uploaded a file "a-folder/a-sub-folder/testfile.txt" with content "text" in the home directory with the alias "testfile.txt"
    Then for user "admin" the treesize of the resource with the alias "Admin Home" should be 4
    When user "admin" has uploaded a file "a-folder/a-sub-folder-B/testfile2.txt" with content "text" in the home directory with the alias "testfile2.txt"
    Then for user "admin" the treesize of the resource with the alias "Admin Home" should be 8
    And for user "admin" the treesize of the resource with the alias "a-folder" should be 8
    And for user "admin" the treesize of the resource with the alias "a-sub-folder" should be 4
    And for user "admin" the treesize of the resource with the alias "a-sub-folder-B" should be 4

 Scenario: Uploading new content should change the checksums
    Given user "admin" has uploaded a file "a-folder/a-sub-folder/testfile.txt" with content "text" in the home directory with the alias "testfile.txt"
    When user "admin" remembers the fileinfo of the resource with the alias "testfile.txt"
    # we leave the alias empty because we do not want to overwrite the remembered file info
    And user "admin" has uploaded a file "a-folder/a-sub-folder/testfile.txt" with content "new text" in the home directory with the alias ""
    Then for user "admin" the checksums of the resource with the alias "testfile.txt" should have changed
