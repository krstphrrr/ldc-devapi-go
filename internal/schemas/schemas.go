package schemas


var TableSchemas = map[string]map[string]string{
	"public_test.datalpi_filtered_view":           DataLPISchema,
	"public_test.dataheader_filtered_view":        DataHeaderSchema,
	"public_test.dataheight_filtered_view":        DataHeightSchema,
	"public_test.datagap_filtered_view":           DataGapSchema,
	"public_test.datahorizontalflux_filtered_view": DataHorizontalFluxSchema,
	"public_test.dataplotcharacterization_filtered_view": DataPlotCharacterizationSchema,
	"public_test.datasoilhorizons_filtered_view":  DataSoilHorizonsSchema,
	"public_test.datasoilstability_filtered_view": DataSoilStabilitySchema,
	"public_test.dataspeciesinventory_filtered_view": DataSpeciesInventorySchema,
	"public_test.geoindicators_filtered_view":     GeoIndicatorsSchema,
	"public_test.geospecies_filtered_view":        GeoSpeciesSchema,
	"public_test.tblrhem_filtered_view":           TblRHEMSchema,
	"public_test.tblproject_filtered_view":        TblProjectSchema,
	"aero_data.aero_summary":                      DataAeroSummarySchema,
}
