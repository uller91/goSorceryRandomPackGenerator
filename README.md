# SorceryPackGen: Sealed card packs generator for Sorcery TCG

Allows you to generate custom card packs for various sets of **Sorcery TCG**

## Motivation

If you, like me, prefer to play limited formats in **Sorcery TCG** but you already have full collection (or just don't want to spend more money on cards packs), **SorceryPackGen** is a CLI-tool that gives you the ability to do so easilty (or with as little effort as possible). This app's aim is to generate random standard or custom card packs for **Sorcery TCG** which are useful in the formats like **sealed** or **draft**.

## Example of use

### Generating one standard pack from Alpha set.

**Input:**
```
sorceryPackGen generate -s Alpha -p standard -n 1
```
**Output:**
```
Random pack from Alpha set:

Highland Clansmen         | Minion     | Ordinary       
Swan Maidens              | Minion     | Ordinary       
Cave Trolls               | Minion     | Ordinary       
Humble Village            | Site       | Ordinary       
Ice Lance                 | Magic      | Ordinary       
Rustic Village            | Site       | Ordinary       
Scourge Zombies           | Minion     | Ordinary       
Summer River              | Site       | Ordinary       
Common Sense              | Magic      | Ordinary       
Sacred Scarabs            | Minion     | Ordinary       
Sandstorm                 | Aura       | Ordinary       
Ormund Harpooneers        | Minion     | Exceptional    
Giant Shark               | Minion     | Exceptional    
Iceberg                   | Site       | Exceptional    
Scorched Earth            | Magic      | Elite    
```

## Prerequisites

1. Go v.1.24.1+. Available [here](https://webinstall.dev/golang/) or [here](https://go.dev/doc/install)
2. PostgreSQL database. Available [here](https://webinstall.dev/postgres/) or [here](https://learn.microsoft.com/en-us/windows/wsl/tutorials/wsl-database#install-postgresql)
3. Goose CLI DB migration tool. Available [here](https://github.com/pressly/goose#install)

## Quick Start

1. Clone the repo:

       git clone https://github.com/uller91/sorceryPackGen.git

2. Set-up **postgreSQL** DB for **sorceryPackGen** cards DB.
3. Copy the content of **.env_example** into a new **.env** file changing the **example connection string** to the one for the **postgreSQL DB** you just set-up.
4. Perform the DB migration in from the same dirrectory using the command below:
        
       goose -dir sql/schema postgres *your_DB_connection_string* up

5. Install the app:

       go install

6. Check that the app was installed properly:

       sorceryPackGen version

7. Update the internal card DB:

       sorceryPackGen update

8. Generate your first card pack:

       sorceryPackGen generate

## Usage

* To see the list of all available commands.
       sorceryPackGen help
* `sorceryPackGen generate` to generate one standard 15 cards pack from random Sorcery TCG set. 
* `sorceryPackGen help generate` to see the list of all available pack generation options.

## Contributing

If you'd like to contribute to this project, feel free to fork the repo and submit pull requests for bug fixes, enhancements, or documentation improvements to the **main** brunch.

## Future plans

Future plans include support of all future Sorcery TCG as well as custom card collections.