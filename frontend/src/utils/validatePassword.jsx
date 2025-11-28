export const validatePassword = (password) => {
    if (password.length < 8) return false;

    const hasUpper = /[A-Z]/.test(password);
    const hasLower = /[a-z]/.test(password);
    const hasDigit = /\d/.test(password);
    const hasSymbol = /[^A-Za-z0-9]/.test(password);

    return hasUpper && hasLower && hasDigit && hasSymbol;
}