# PDF Generation Documentation

This directory contains the PDF generation components for form printing.

## Architecture

### Files
- **`FormPDFDocument.tsx`** - The main PDF document component that defines the layout and content of the generated PDF

### How It Works

1. User clicks "Print Format" button on forms list page
2. Navigates to `/forms/[id]/print` route
3. Print page fetches form data from API
4. Data is processed/transformed as needed (TypeScript logic)
5. Processed data is passed to `FormPDFDocument` component
6. PDF can be previewed in browser or downloaded

## Adding Custom TypeScript Logic

You can add any data processing logic in two places:

### 1. In the Print Page (`/app/forms/[id]/print/page.tsx`)

```typescript
// Process data before passing to PDF component
const processedData = form ? {
    ...form,
    // Add computed fields
    formattedCreatedAt: new Date(form.created_at).toLocaleDateString(),
    fullName: `${form.first_name} ${form.last_name}`,
    // Calculate aggregates
    totalAmount: form.pest_apps.reduce((sum, app) => sum + parseFloat(app.amount_applied), 0),
    // Filter or transform arrays
    recentApps: form.pest_apps.filter(app => isRecent(app.app_timestamp)),
    // Any other TypeScript logic
} : null;
```

### 2. In the PDF Document Component (`FormPDFDocument.tsx`)

```typescript
const FormPDFDocument: React.FC<FormPDFDocumentProps> = ({ form }) => {
    // Add calculations
    const totalApplications = form.pest_apps?.length || 0;

    // Group data
    const appsByLocation = form.pest_apps.reduce((acc, app) => {
        acc[app.location_code] = [...(acc[app.location_code] || []), app];
        return acc;
    }, {});

    // Format data
    const displayDate = new Date(form.created_at).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });

    // Continue with PDF rendering...
}
```

## Styling the PDF

Styles are defined using `StyleSheet.create()` from `@react-pdf/renderer`:

```typescript
const styles = StyleSheet.create({
    page: {
        padding: 30,
        fontSize: 12,
        fontFamily: 'Helvetica',
    },
    header: {
        fontSize: 20,
        marginBottom: 20,
        textAlign: 'center',
        fontWeight: 'bold',
    },
    // Add more styles...
});
```

## Available PDF Components

From `@react-pdf/renderer`:
- `<Document>` - Root component
- `<Page>` - Individual pages
- `<View>` - Container (like div)
- `<Text>` - Text content
- `<Image>` - Images (requires src URL or base64)
- `<Link>` - Clickable links

## Adding Sections

To add a new section to the PDF:

```typescript
<View style={styles.section}>
    <Text style={styles.sectionTitle}>New Section Title</Text>

    <View style={styles.row}>
        <Text style={styles.label}>Label:</Text>
        <Text style={styles.value}>Value</Text>
    </View>

    {/* Add more content */}
</View>
```

## Tables

Create tables using flexbox layout:

```typescript
{/* Header */}
<View style={[styles.tableRow, styles.tableHeader]}>
    <Text style={styles.tableCell}>Column 1</Text>
    <Text style={styles.tableCell}>Column 2</Text>
</View>

{/* Rows */}
{data.map((item, index) => (
    <View key={index} style={styles.tableRow}>
        <Text style={styles.tableCell}>{item.field1}</Text>
        <Text style={styles.tableCell}>{item.field2}</Text>
    </View>
))}
```

## Conditional Rendering

Use standard JavaScript conditionals:

```typescript
{form.form_type === 'lawn' && (
    <View style={styles.row}>
        <Text style={styles.label}>Lawn Area:</Text>
        <Text style={styles.value}>{form.lawn_area_sq_ft} sq ft</Text>
    </View>
)}
```

## Tips

- **Pagination**: Content automatically flows to new pages when needed
- **Fonts**: Default fonts are Helvetica, Times, and Courier. Custom fonts can be registered
- **Images**: Must be URLs or base64 encoded
- **Styling**: Uses CSS-like syntax but not all CSS properties are supported
- **Testing**: Use the preview mode to see changes before downloading

## Common Customizations

### Add Company Logo
```typescript
import { Image } from '@react-pdf/renderer';

<View style={styles.header}>
    <Image src="/path/to/logo.png" style={{ width: 100 }} />
    <Text>Company Name</Text>
</View>
```

### Add Signature Lines
```typescript
<View style={styles.signatureSection}>
    <View style={styles.signatureLine}>
        <Text>Customer Signature: _____________________</Text>
        <Text>Date: _____________</Text>
    </View>
</View>
```

### Multi-Page Documents
```typescript
<Document>
    <Page size="A4" style={styles.page}>
        {/* Page 1 content */}
    </Page>
    <Page size="A4" style={styles.page}>
        {/* Page 2 content */}
    </Page>
</Document>
```

## Resources

- [@react-pdf/renderer Documentation](https://react-pdf.org/)
- [Styling Guide](https://react-pdf.org/styling)
- [Components API](https://react-pdf.org/components)
