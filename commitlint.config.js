// This config should help enforce Chris Bean's commit message recommendations https://cbea.ms/git-commit/
// 
// - Separate subject from body with a blank line
// - Limit the subject line to 50 characters
// - Capitalize the subject line
// - Do not end the subject line with a period
// - Use the imperative mood in the subject line
// - Wrap the body at 72 characters
// - Use the body to explain what and why vs. how

module.exports = {
    extends: ['@commitlint/config-conventional'],
    rules: {
        'body-leading-blank': [2, 'always'],
        'header-max-length': [2, 'always', 50],
        'header-capitalized': [2, 'always'],
        'header-full-stop': [2, 'never'],
        'header-case': [2, 'always', 'sentence-case'],
        'subject-case': [2, 'always', 'sentence-case'],
        'subject-full-stop': [0, 'never'],
        'imperative-mood': [2, 'always'],
    },
    parserPreset: {
        parserOpts: {
            headerPattern: /^(\w*)(?:\((.*)\))?-(.*)$/,
            headerCorrespondence: ['type', 'scope', 'subject'],
        },
    },
    plugins: [
        {
            rules: {
                'imperative-mood': ({ subject }) => {
                    const imperativeMoodRegex = /^(add|create|update|fix|remove|refactor|improve|...) .+$/i;
                    return imperativeMoodRegex.test(subject)
                        ? [true]
                        : [`Subject must use the imperative mood.`];
                },
            },
        },
    ],
};
