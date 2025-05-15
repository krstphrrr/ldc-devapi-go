package schemas

var TableSchemas = map[string]map[string]string{
	"public_dev.datalpi_filtered_view":                  DataLPISchema,
	"public_dev.dataheader_filtered_view":               DataHeaderSchema,
	"public_dev.dataheight_filtered_view":               DataHeightSchema,
	"public_dev.datagap_filtered_view":                  DataGapSchema,
	"public_dev.datahorizontalflux_filtered_view":       DataHorizontalFluxSchema,
	"public_dev.dataplotcharacterization_filtered_view": DataPlotCharacterizationSchema,
	"public_dev.datasoilhorizons_filtered_view":         DataSoilHorizonsSchema,
	"public_dev.datasoilstability_filtered_view":        DataSoilStabilitySchema,
	"public_dev.dataspeciesinventory_filtered_view":     DataSpeciesInventorySchema,
	"public_dev.geoindicators_filtered_view":            GeoIndicatorsSchema,
	"public_dev.geospecies_filtered_view":               GeoSpeciesSchema,
	"public_dev.tblrhem_filtered_view":                  TblRHEMSchema,
	"public_dev.tblproject_filtered_view":               TblProjectSchema,
	"aero_data.aero_summary":                            DataAeroSummarySchema,
}
