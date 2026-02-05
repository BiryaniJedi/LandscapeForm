/**
 * Site location codes for pesticide application forms.
 *
 * Codes follow a two-part system:
 * - Numeric prefix (1-5): Property area identifier
 * - Alphabetic suffix (A-D): Feature type identifier
 *
 * Example: "1A" = Front Yard, Lawn
 */

export const siteCodesFirst = {
    '1': 'Front Yard',
    '2': 'Side Yard(s)',
    '3': 'Rear Yard',
    '4': 'Entire Property',
    '5': 'Other',
}
export const siteCodesSecond = {
    'A': 'Lawn',
    'B': 'Shrub(s)',
    'C': 'Ornamental Tree(s)',
    'D': 'Other',
}
