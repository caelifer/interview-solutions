#/usr/bin/perl -w

# Goat Latin is a made-up language based off of English, sort of like Pig Latin.
# The rules of Goat Latin are as follows:
# 1. If a word begins with a consonant (i.e. not a vowel), remove the first
#    letter and append it to the end, then add 'ma'.
#    For example, the word 'goat' becomes 'oatgma'.
# 2. If a word begins with a vowel, append 'ma' to the end of the word.
#    For example, the word 'I' becomes 'Ima'.
# 3. Add one letter "a" to the end of each word per its word index in the
#    sentence, starting with 1. That is, the first word gets "a" added to the
#    end, the second word gets "aa" added to the end, the third word in the
#    sentence gets "aaa" added to the end, and so on.

# Write a function that, given a string of words making up one sentence, returns
# that sentence in Goat Latin. For example:
#
#  string_to_goat_latin('I speak Goat Latin')
#
# would return: 'Imaa peaksmaaa oatGmaaaa atinLmaaaaa'

use strict;

my $seps = qr{[ \t\n]}o;

sub isVowel {
    shift() =~ /[aeiou]/i ? 1 : 0;
}

sub string_to_goat_latin {
    my @words = split(/$seps/, shift());
    
    for (my $i = 0; $i < scalar(@words);) {
        # Get first char
        my @letters = split('', $words[$i]);
        
        if (! isVowel($letters[0])) {
            # Rule #1 - consonants - rotate first char
            my $fc = shift @letters;
            push @letters, $fc;
        }
        
        # Rule #2 and second part of Rule #1 add 'ma' to the end
        push(@letters, 'm', 'a');
        
        # Rule #3 - Add 'a' x word index
        print(join('', @letters) . 'a' x ++$i);

        # Print word separator
        print " " if $i < scalar(@words);
    }

    print "\n";
}

string_to_goat_latin('I speak Goat Latin');

