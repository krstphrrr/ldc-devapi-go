package schemas

// dataSoilStabilitySchema defines the column names and types for the `dataSoilStability` table.
var DataSoilStabilitySchema = map[string]string{
    "rid": "integer",
    "PrimaryKey": "text",
    "DBKey": "text",
    "ProjectKey": "text",
    "RecKey": "text",
    "FormType": "text",
    "FormDate": "date",
    "LineKey": "text",
    "SoilStabSubSurface": "integer",
    "Line": "text",
    "Position": "integer",
    "Pos": "text",
    "Veg": "text",
    "Rating": "integer",
    "Hydro": "integer",
    "Notes": "text",
    "source": "text",
    "DateLoadedInDb": "date",
    "DateVisited": "date",
}
