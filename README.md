# FURY DB

Pure GO lang embedded SQL database made quickly and furiously. For hacks and fun.

## Background

On the 29th Dec 2020 to the 31th. For 3 days at the Ninja Software we had a Hackathon and this is my project.

## Usage

Test codes in cmd dir

## Releases

After the demo and at the moment, I'm planning to slowly to implement features for learning and fun, no guarantee though.

## Features

I would like FuryDB to have at least following features (subject to change), so it can at least be used in basic projects.

x = done  
p = partially done

- Table
    - [p] Create
    - [ ] Alter
    - [ ] Delete
    - [ ] Types
        - [x] Boolean
        - [x] Int
        - [x] Float
        - [x] String
        - [ ] Time
        - [p] UUID
    - Constraints
        - [p] Primary
        - [ ] Foreign key
        - [ ] Nullable
        - [ ] Default
        - [ ] Unique
- Record
    - [p] Insert
    - [ ] Update
    - [ ] Delete
    - [p] Select
    - Condition
        - [ ] Where
        - [ ] Limit
        - [ ] Top
        - [ ] Offset
        - [ ] Distinct
        - [ ] Order By
        - [ ] And
        - [ ] Or
        - [ ] In
        - [ ] IS NULL
        - [ ] Comparison
            - [ ] Equal
            - [ ] Larger (and equal) than
            - [ ] Less (and equal) than
            - [ ] Like
        - [ ] Group by
    - [ ] Left Join other table
    - Aggregate functions
        - [ ] Count
        - [ ] Sum
        - [ ] Min
        - [ ] Max
        - [ ] Sum
- Go SQL Driver
    - [x] Open
    - [ ] Close
    - Query
        - [p] Basic
        - [ ] parameterized query
    - [ ] Exec
        - [p] Basic
        - [ ] parameterized query
