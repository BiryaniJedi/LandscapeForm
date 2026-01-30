export const formatDate = (dateString: string) => {
    const dateObject = new Date(dateString);
    return dateObject.toLocaleString();
}
