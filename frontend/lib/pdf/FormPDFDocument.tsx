import React from 'react';
import { Document, Page, Text, View, StyleSheet } from '@react-pdf/renderer';
import { FormViewResponse, ListChemicalsResponse } from '@/lib/api/types';
import { siteCodesFirst, siteCodesSecond } from '../common/siteCodes';
import { chemicalsClient } from '../api/chemicals';

// Compact styles for single-page layout
const styles = StyleSheet.create({
    page: {
        padding: 15,
        fontSize: 8,
        fontFamily: 'Helvetica',
    },
    titleBox: {
        backgroundColor: '#FFFFFF',
        border: '2px solid #000000',
        padding: 12,
        marginBottom: 8,
        minHeight: 35,
        width: '50%',
    },
    titleText: {
        fontSize: 12,
        fontWeight: 'bold',
        textAlign: 'center',
    },
    row: {
        flexDirection: 'row',
        marginBottom: 8,
        gap: 8,
    },
    box: {
        backgroundColor: '#FFFFFF',
        border: '1px solid #000000',
        padding: 6,
    },
    boxTitle: {
        fontSize: 10,
        fontWeight: 'bold',
        marginBottom: 4,
        borderBottom: '1px solid #000',
        paddingBottom: 2,
    },
    field: {
        marginBottom: 3,
    },
    fieldValue: {
        backgroundColor: '#FFFFFF',
        border: '1px solid #000000',
        padding: 3,
        minHeight: 15,
    },
    fieldLabel: {
        fontSize: 7,
        color: '#000000',
        marginTop: 1,
    },
    table: {
        marginTop: 4,
    },
    tableRow: {
        flexDirection: 'row',
        borderBottom: '1px solid #000',
        paddingVertical: 2,
    },
    tableHeader: {
        backgroundColor: '#e0e0e0',
        fontWeight: 'bold',
        fontSize: 7,
    },
    tableCell: {
        flex: 1,
        paddingHorizontal: 2,
        fontSize: 7,
    },
    smallBox: {
        backgroundColor: '#FFFFFF',
        border: '1px solid #000000',
        padding: 6,
        minHeight: 60,
    },
    emptyBoxLabel: {
        fontSize: 8,
        color: '#666',
        fontStyle: 'italic',
    },
    codeTable: {
        marginTop: 2,
    },
    codeRow: {
        flexDirection: 'row',
        borderBottom: '1px solid #ccc',
        paddingVertical: 2,
    },
    codeCell: {
        fontSize: 7,
        paddingHorizontal: 2,
    },
    codeNumber: {
        width: 20,
        fontWeight: 'bold',
    },
    codeDesc: {
        flex: 1,
    },
    sideBoxesRow: {
        flexDirection: 'row',
        gap: 4,
        marginTop: 2,
    },
    sideBox: {
        flex: 1,
    },
});

interface FormPDFDocumentProps {
    form: FormViewResponse & {
        formattedCreatedAt?: string;
        formattedUpdatedAt?: string;
        fullName?: string;
        fullAddress?: string;
    };
    chemicalList: ListChemicalsResponse;
}

const FormPDFDocument: React.FC<FormPDFDocumentProps> = ({ form, chemicalList }) => {
    // Get unique chemicals used in this form
    return (
        <Document>
            <Page size="A4" style={styles.page}>
                {/* Row 1: Title + Customer Info */}
                <View style={styles.row}>
                    {/* Title Box - 50% */}
                    <View style={[styles.titleBox, { width: '60%' }]}>
                        <Text style={styles.titleText}>
                            {form.form_type.toUpperCase()} Treatment Program April - November 2026
                        </Text>
                        <Text style={{ fontSize: 8, textAlign: 'center', marginTop: 4 }}>
                            Kindergan Address Here
                        </Text>
                        <Text style={{ fontSize: 8, textAlign: 'center', marginTop: 2 }}>
                            Kindergan Address Here
                        </Text>
                        <Text style={{ fontSize: 8, textAlign: 'center', marginTop: 2 }}>
                            Kindergan Address Here
                        </Text>
                        <Text style={{ fontSize: 8, textAlign: 'center', marginTop: 2 }}>
                            NJDEP#: **number**
                        </Text>
                        <Text style={{ fontSize: 8, textAlign: 'center', marginTop: 2 }}>
                            Resp Applic: **Name** - Lic #: **number**
                        </Text>
                        <Text style={{ fontSize: 8, textAlign: 'center', marginTop: 2 }}>
                            Operator: **Name** - Lic #: **number**
                        </Text>
                    </View>
                    {/* Customer Information - 50% */}
                    <View style={[styles.box, { width: '40%' }]}>
                        <Text style={styles.boxTitle}>Customer Information</Text>
                        <View style={styles.field}>
                            <Text style={styles.fieldValue}>
                                {form.fullName || `${form.first_name} ${form.last_name}`}
                            </Text>
                            <Text style={styles.fieldLabel}>Name</Text>
                        </View>
                        <View style={styles.field}>
                            <Text style={styles.fieldValue}>
                                {form.fullAddress || `${form.street_number} ${form.street_name}, ${form.town} ${form.zip_code}`}
                            </Text>
                            <Text style={styles.fieldLabel}>Address</Text>
                        </View>
                        <View style={{ flexDirection: 'row', gap: 4 }}>
                            <View style={[styles.field, { flex: 1 }]}>
                                <Text style={styles.fieldValue}>{form.home_phone}</Text>
                                <Text style={styles.fieldLabel}>Home Phone</Text>
                            </View>
                            {form.other_phone && (
                                <View style={[styles.field, { flex: 1 }]}>
                                    <Text style={styles.fieldValue}>{form.other_phone}</Text>
                                    <Text style={styles.fieldLabel}>Other Phone</Text>
                                </View>
                            )}
                        </View>
                        <View style={{ flexDirection: 'row', gap: 4 }}>
                            <View style={[styles.field, { flex: 1 }]}>
                                <Text style={styles.fieldValue}>{form.call_before ? 'Yes' : 'No'}</Text>
                                <Text style={styles.fieldLabel}>Call Before</Text>
                            </View>
                            <View style={[styles.field, { flex: 1 }]}>
                                <Text style={styles.fieldValue}>{form.is_holiday ? 'Yes' : 'No'}</Text>
                                <Text style={styles.fieldLabel}>Jewish Holiday</Text>
                            </View>
                            {form.form_type === 'lawn' && form.lawn_area_sq_ft !== undefined && (
                                <>
                                    <View style={[styles.field, { flex: 1 }]}>
                                        <Text style={styles.fieldValue}>{form.lawn_area_sq_ft} sq ft</Text>
                                        <Text style={styles.fieldLabel}>Lawn Area</Text>
                                    </View>
                                    <View style={[styles.field, { flex: 1 }]}>
                                        <Text style={styles.fieldValue}>{form.fert_only ? 'Yes' : 'No'}</Text>
                                        <Text style={styles.fieldLabel}>Fertilizer Only</Text>
                                    </View>
                                </>
                            )}
                            {form.form_type === 'shrub' && form.flea_only !== undefined && (
                                <View style={[styles.field, { flex: 1 }]}>
                                    <Text style={styles.fieldValue}>{form.flea_only ? 'Yes' : 'No'}</Text>
                                    <Text style={styles.fieldLabel}>Flea Only</Text>
                                </View>
                            )}
                        </View>
                    </View>
                </View>
                {/* Row 2: Chemical List and Site Codes */}
                <View style={styles.row}>
                    {/* Chemical List - 65% */}
                    <View style={[styles.box, { width: '65%' }]}>
                        <Text style={styles.boxTitle}>Chemicals</Text>
                        <View style={styles.table}>
                            <View style={[styles.tableRow, styles.tableHeader]}>
                                <Text style={[styles.tableCell, { flex: 1.5 }]}>ID</Text>
                                <Text style={[styles.tableCell, { flex: 1.5 }]}>Brand Name</Text>
                                <Text style={[styles.tableCell, { flex: 1.5 }]}>Chemical Name</Text>
                                <Text style={[styles.tableCell, { flex: 1.5 }]}>EPA Reg. No</Text>
                                <Text style={[styles.tableCell, { flex: 1.5 }]}>Recipe</Text>
                            </View>
                            {chemicalList.chemicals.map((chem, index) => (
                                <View key={chem.id || index} style={styles.tableRow}>
                                    <Text style={[styles.tableCell, { flex: 1.5 }]}>{chem.id}</Text>
                                    <Text style={[styles.tableCell, { flex: 1.5 }]}>{chem.brand_name}</Text>
                                    <Text style={[styles.tableCell, { flex: 1.5 }]}>{chem.chemical_name}</Text>
                                    <Text style={[styles.tableCell, { flex: 1.5 }]}>{chem.epa_reg_no}</Text>
                                    <Text style={[styles.tableCell, { flex: 1.5 }]}>{chem.recipe}</Text>
                                </View>
                            ))}
                        </View>
                    </View>
                    {/* Site Codes - 35% */}
                    <View style={[styles.box, { width: '35%' }]}>
                        <Text style={styles.boxTitle}>Site Code Reference</Text>
                        <View style={styles.sideBoxesRow}>
                            {/* First Codes */}
                            <View style={styles.sideBox}>
                                <Text style={{ fontSize: 8, fontWeight: 'bold', marginBottom: 2 }}>First Code</Text>
                                {Object.entries(siteCodesFirst).map(([code, desc]) => (
                                    <View key={code} style={styles.codeRow}>
                                        <Text style={[styles.codeCell, styles.codeNumber]}>{code}</Text>
                                        <Text style={[styles.codeCell, styles.codeDesc]}>{desc}</Text>
                                    </View>
                                ))}
                            </View>
                            {/* Second Codes */}
                            <View style={styles.sideBox}>
                                <Text style={{ fontSize: 8, fontWeight: 'bold', marginBottom: 2 }}>Second Code</Text>
                                {Object.entries(siteCodesSecond).map(([code, desc]) => (
                                    <View key={code} style={styles.codeRow}>
                                        <Text style={[styles.codeCell, styles.codeNumber]}>{code}</Text>
                                        <Text style={[styles.codeCell, styles.codeDesc]}>{desc}</Text>
                                    </View>
                                ))}
                            </View>
                        </View>
                    </View>
                </View>

                <View style={styles.box}>
                    <Text style={styles.boxTitle}>Pesticide Applications</Text>
                    {form.pest_apps && form.pest_apps.length > 0 ? (
                        <View style={styles.table}>
                            <View style={[styles.tableRow, styles.tableHeader]}>
                                <Text style={[styles.tableCell, { flex: 1.5 }]}>Date</Text>
                                <Text style={[styles.tableCell, { flex: 1.5 }]}>Chemical</Text>
                                <Text style={[styles.tableCell, { flex: 1 }]}>Rate</Text>
                                <Text style={[styles.tableCell, { flex: 1 }]}>Amount</Text>
                                <Text style={[styles.tableCell, { flex: 1 }]}>Location</Text>
                            </View>
                            {form.pest_apps.map((app, index) => (
                                <View key={app.id || index} style={styles.tableRow}>
                                    <Text style={[styles.tableCell, { flex: 1.5 }]}>
                                        {new Date(app.app_timestamp).toLocaleDateString()}
                                    </Text>
                                    <Text style={[styles.tableCell, { flex: 1.5 }]}>{app.chem_used}</Text>
                                    <Text style={[styles.tableCell, { flex: 1 }]}>{app.rate}</Text>
                                    <Text style={[styles.tableCell, { flex: 1 }]}>{app.amount_applied.toString()}</Text>
                                    <Text style={[styles.tableCell, { flex: 1 }]}>{app.location_code}</Text>
                                </View>
                            ))}
                        </View>
                    ) : (
                        <Text style={{ fontSize: 7, fontStyle: 'italic' }}>No applications recorded</Text>
                    )}
                </View>

                {/* Row 3: Chemical List + Site Codes */}
                <View style={styles.row}>
                </View>

                {/* Row 4: Two Empty Boxes */}
                <View style={styles.row}>
                    <View style={[styles.box, { flex: 1 }]}>
                        <Text style={{ fontSize: 10, fontWeight: 'bold', marginTop: 2, textAlign: 'center' }}>
                            New Jersey D.E.P. Pesticide Control Program: (609) 984-6507
                        </Text>
                        <Text style={{ fontSize: 10, fontWeight: 'bold', marginTop: 2, textAlign: 'center' }}>
                            National Pesticide Information Center: (800) 858-7378 (General Questions)
                        </Text>
                        <Text style={{ fontSize: 10, fontWeight: 'bold', marginTop: 2, textAlign: 'center' }}>
                            National Pesticide Information & Education System: (800) 222-1222 (Emergencies)
                        </Text>
                    </View>
                </View>
            </Page>
        </Document>
    );
};

export default FormPDFDocument;
