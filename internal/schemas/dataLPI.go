package schemas

// dataLPISchema defines the column names and types for the `dataLPI` table.
var DataLPISchema = map[string]string{
    "rid": "integer",
    "PrimaryKey": "text",
    "DBKey": "text",
    "ProjectKey": "text",
    "LineKey": "text",
    "RecKey": "text",
    "layer": "text",
    "code": "text",
    "chckbox": "integer",
    "ShrubShape": "text",
    "FormType": "text",
    "FormDate": "date",
    "Direction": "text",
    "Measure": "integer",
    "LineLengthAmount": "integer",
    "SpacingIntervalAmount": "double precision",
    "SpacingType": "text",
    "ShowCheckbox": "integer",
    "CheckboxLabel": "text",
    "PointLoc": "double precision",
    "PointNbr": "integer",
    "source": "text",
    "DateLoadedInDb": "date",
    "DateVisited": "date",
}
