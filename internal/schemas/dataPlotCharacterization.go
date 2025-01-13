package schemas

// dataPlotCharacterizationSchema defines the column names and types for the `dataPlotCharacterization` table.
var DataPlotCharacterizationSchema = map[string]string{
    "rid": "integer",
    "ProjectKey": "text",
    "PrimaryKey": "text",
    "DBKey": "text",
    "EstablishDate": "date",
    "State": "text",
    "County": "text",
    "EcolSite": "text",
    "ParentMaterial": "text",
    "Slope": "double precision",
    "Aspect": "double precision",
    "LandscapeType": "text",
    "LandscapeTypeSecondary": "text",
    "Elevation": "double precision",
    "SoilSeries": "text",
    "Longitude_NAD83": "double precision",
    "Latitude_NAD83": "double precision",
    "SlopeShapeVertical": "text",
    "SlopeShapeHorizontal": "text",
    "MLRA": "text",
    "source": "text",
    "DateVisited": "date",
    "DateLoadedInDb": "date",
}
