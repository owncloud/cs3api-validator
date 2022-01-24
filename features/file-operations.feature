Feature: Basic File operations
  As a user
  I want to be able to do basic file operations
  So that I can organize my files

  Background: Space Exists
    Given user "marie" has logged in with password "radioactivity"
    Given user "marie" has created a personal space with the alias "Maries Home"

    Scenario: Create Folders and files
      Given user "marie" has created a folder "NewFolder" in the home directory with the alias "NewFolder"
      And user "marie" has uploaded a file "testfile.txt" with content "lorem ipsum" in the home directory with the alias "testfile"
      When user "marie" lists all resources inside the resource with alias "Maries Home"
      Then the following resources should be listed in the response:
        | type      | path         |
        | container | NewFolder    |
        | file      | testfile.txt |

    Scenario: Create Tree
      Given user "marie" has created a folder "A" in the home directory with the alias "A"
      And user "marie" has created a folder "A/B" in the home directory with the alias "B"
      And user "marie" has uploaded a file "A/B/testfile.txt" with content "lorem ipsum" in the home directory with the alias "testfile"
      When user "marie" lists all resources inside the resource with alias "Maries Home"
      Then the following resources should not be listed in the response:
        | type      | path         |
        | container | B            |
        | file      | testfile.txt |
      And the following resources should be listed in the response:
        | type      | path |
        | container | A    |

  Scenario: Rename resources
    Given user "marie" has created a folder "A" in the home directory with the alias "A"
    And user "marie" has uploaded a file "testfile.txt" with content "lorem ipsum" in the home directory with the alias "testfile"
    When user "marie" lists all resources inside the resource with alias "Maries Home"
    Then the following resources should be listed in the response:
      | type      | path         |
      | container | A            |
      | file      | testfile.txt |
    When user "marie" has moved the resource with alias "A" inside a space to target "A-Folder"
    And user "marie" has moved the resource with alias "testfile" inside a space to target "Test-File.txt"
    When user "marie" lists all resources inside the resource with alias "Maries Home"
    Then the following resources should be listed in the response:
      | type      | path          |
      | container | A-Folder      |
      | file      | Test-File.txt |
    And the following resources should not be listed in the response:
      | type      | path         |
      | container | A            |
      | file      | testfile.txt |
