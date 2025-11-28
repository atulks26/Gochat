import bcrypt from 'bcryptjs';

export async function passwordHash (password) {
    const saltRounds = 10;
    const salt = await bcrypt.genSalt(saltRounds);
    const hashed = await bcrypt.hash(password, salt);

    return hashed;
}