package schemas

// geoSpeciesSchema defines the column names and types for the `geoSpecies` table.
var GeoSpeciesSchema = map[string]string{
    "rid": "integer",
    "PrimaryKey": "text",
    "DBKey": "text",
    "ProjectKey": "text",
    "Species": "text",
    "AH_SpeciesCover": "double precision",
    "AH_SpeciesCover_n": "integer",
    "Hgt_Species_Avg": "double precision",
    "Hgt_Species_Avg_n": "integer",
    "Duration": "text",
    "GrowthHabit": "text",
    "GrowthHabitSub": "text",
    "SpeciesKey": "text",
    "DateLoadedInDb": "date",
    "DateVisited": "date",
    "ScientificName": "text",
}
