package schemas

// dataHeaderSchema defines the column names and types for the `dataHeader` table.
var DataHeaderSchema = map[string]string{
    "rid": "integer",
    "PrimaryKey": "text",
    "DBKey": "text",
    "ProjectKey": "text",
    "DateVisited": "date",
    "Latitude_NAD83": "double precision",
    "Longitude_NAD83": "double precision",
    "LocationType": "text",
    "EcologicalSiteID": "text",
    "PercentCoveredByEcoSite": "double precision",
    "SpeciesKey": "text",
    "PlotID": "text",
    "DateLoadedInDb": "date",
    "wkb_geometry": "text",
    "source": "text",
}
